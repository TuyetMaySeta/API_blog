package services

import (
	"blog-api/internal/database"
	"blog-api/internal/models"
	"encoding/json"
	"fmt"
	"log"
)

type SearchService struct {
}

func NewSearchService() *SearchService {
	return &SearchService{}
}

// IndexPost indexes a post in Elasticsearch
func (ss *SearchService) IndexPost(post *models.Post) error {
	esPost := models.ElasticsearchPost{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
		Tags:    post.Tags,
	}

	err := database.IndexPost(esPost, fmt.Sprintf("%d", post.ID))
	if err != nil {
		log.Printf("Error indexing post %d: %v", post.ID, err)
		return err
	}

	log.Printf("Post %d indexed successfully", post.ID)
	return nil
}

// DeletePost removes a post from Elasticsearch index
func (ss *SearchService) DeletePost(id uint) error {
	err := database.DeletePost(fmt.Sprintf("%d", id))
	if err != nil {
		log.Printf("Error deleting post %d from index: %v", id, err)
		return err
	}

	log.Printf("Post %d deleted from index successfully", id)
	return nil
}

// SearchPosts performs full-text search on posts
func (ss *SearchService) SearchPosts(query string) ([]models.Post, error) {
	searchResult, err := database.SearchPosts(query)
	if err != nil {
		log.Printf("Error searching posts: %v", err)
		return nil, err
	}

	var posts []models.Post
	for _, hit := range searchResult.Hits.Hits {
		var esPost models.ElasticsearchPost
		err := json.Unmarshal(hit.Source, &esPost)
		if err != nil {
			log.Printf("Error unmarshaling search result: %v", err)
			continue
		}

		post := models.Post{
			ID:      esPost.ID,
			Title:   esPost.Title,
			Content: esPost.Content,
			Tags:    esPost.Tags,
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// FindRelatedPosts finds posts with similar tags
func (ss *SearchService) FindRelatedPosts(tags []string, excludeID uint, limit int) ([]models.Post, error) {
	if len(tags) == 0 {
		return []models.Post{}, nil
	}

	searchResult, err := database.FindRelatedPosts(tags, excludeID, limit)
	if err != nil {
		log.Printf("Error finding related posts: %v", err)
		return nil, err
	}

	var posts []models.Post
	for _, hit := range searchResult.Hits.Hits {
		var esPost models.ElasticsearchPost
		err := json.Unmarshal(hit.Source, &esPost)
		if err != nil {
			log.Printf("Error unmarshaling related post: %v", err)
			continue
		}

		post := models.Post{
			ID:      esPost.ID,
			Title:   esPost.Title,
			Content: esPost.Content,
			Tags:    esPost.Tags,
		}
		posts = append(posts, post)
	}

	return posts, nil
}