package services

import (
	"blog-api/internal/database"
	"blog-api/internal/models"
	"log"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type PostService struct {
	db           *gorm.DB
	cacheService *CacheService
	searchService *SearchService
}

func NewPostService() *PostService {
	return &PostService{
		db:            database.GetDB(),
		cacheService:  NewCacheService(),
		searchService: NewSearchService(),
	}
}

// CreatePost creates a new post with transaction for data integrity
func (ps *PostService) CreatePost(req *models.CreatePostRequest) (*models.Post, error) {
	// Start transaction
	tx := ps.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create post
	post := &models.Post{
		Title:   req.Title,
		Content: req.Content,
		Tags:    pq.StringArray(req.Tags),
	}

	if err := tx.Create(post).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create activity log
	activityLog := &models.ActivityLog{
		Action: "new_post",
		PostID: post.ID,
	}

	if err := tx.Create(activityLog).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Index post in Elasticsearch (async - don't fail if this fails)
	go func() {
		if err := ps.searchService.IndexPost(post); err != nil {
			log.Printf("Failed to index post %d: %v", post.ID, err)
		}
	}()

	return post, nil
}

// GetPostByID retrieves a post by ID with Cache-Aside pattern
func (ps *PostService) GetPostByID(id uint) (*models.PostResponse, error) {
	// Try to get from cache first
	cachedPost, err := ps.cacheService.GetPost(id)
	if err != nil {
		log.Printf("Error getting post from cache: %v", err)
	}

	if cachedPost != nil {
		log.Printf("Cache hit for post %d", id)
		response := &models.PostResponse{
			ID:        cachedPost.ID,
			Title:     cachedPost.Title,
			Content:   cachedPost.Content,
			Tags:      cachedPost.Tags,
			CreatedAt: cachedPost.CreatedAt,
			UpdatedAt: cachedPost.UpdatedAt,
		}

		// Get related posts (bonus feature)
		if len(cachedPost.Tags) > 0 {
			relatedPosts, _ := ps.searchService.FindRelatedPosts(cachedPost.Tags, id, 5)
			response.RelatedPosts = relatedPosts
		}

		return response, nil
	}

	log.Printf("Cache miss for post %d", id)

	// Get from database
	var post models.Post
	if err := ps.db.First(&post, id).Error; err != nil {
		return nil, err
	}

	// Store in cache
	go func() {
		if err := ps.cacheService.SetPost(&post); err != nil {
			log.Printf("Error storing post in cache: %v", err)
		}
	}()

	response := &models.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		Tags:      post.Tags,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}

	// Get related posts (bonus feature)
	if len(post.Tags) > 0 {
		relatedPosts, _ := ps.searchService.FindRelatedPosts(post.Tags, id, 5)
		response.RelatedPosts = relatedPosts
	}

	return response, nil
}

// UpdatePost updates a post and handles cache invalidation
func (ps *PostService) UpdatePost(id uint, req *models.UpdatePostRequest) (*models.Post, error) {
	var post models.Post
	if err := ps.db.First(&post, id).Error; err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}
	if req.Tags != nil {
		post.Tags = pq.StringArray(req.Tags)
	}

	// Save to database
	if err := ps.db.Save(&post).Error; err != nil {
		return nil, err
	}

	// Invalidate cache
	go func() {
		if err := ps.cacheService.InvalidatePost(id); err != nil {
			log.Printf("Error invalidating cache for post %d: %v", id, err)
		}
	}()

	// Update Elasticsearch index
	go func() {
		if err := ps.searchService.IndexPost(&post); err != nil {
			log.Printf("Failed to update post %d in search index: %v", post.ID, err)
		}
	}()

	return &post, nil
}

// DeletePost deletes a post and cleans up cache and search index
func (ps *PostService) DeletePost(id uint) error {
	if err := ps.db.Delete(&models.Post{}, id).Error; err != nil {
		return err
	}

	// Invalidate cache
	go func() {
		if err := ps.cacheService.InvalidatePost(id); err != nil {
			log.Printf("Error invalidating cache for post %d: %v", id, err)
		}
	}()

	// Remove from search index
	go func() {
		if err := ps.searchService.DeletePost(id); err != nil {
			log.Printf("Failed to delete post %d from search index: %v", id, err)
		}
	}()

	return nil
}

// SearchPostsByTag searches posts by tag using GIN index
func (ps *PostService) SearchPostsByTag(tag string) ([]models.Post, error) {
	var posts []models.Post
	
	// Use PostgreSQL's array contains operator with GIN index
	if err := ps.db.Where("tags @> ?", pq.Array([]string{tag})).Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

// SearchPosts performs full-text search using Elasticsearch
func (ps *PostService) SearchPosts(query string) (*models.SearchResponse, error) {
	posts, err := ps.searchService.SearchPosts(query)
	if err != nil {
		return nil, err
	}

	response := &models.SearchResponse{
		Posts: posts,
		Total: int64(len(posts)),
	}

	return response, nil
}

// GetAllPosts retrieves all posts with pagination
func (ps *PostService) GetAllPosts(limit, offset int) ([]models.Post, error) {
	var posts []models.Post
	if err := ps.db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}