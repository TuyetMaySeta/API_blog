package services

import (
	"blog-api/internal/database"
	"blog-api/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheService struct {
	redis *redis.Client
}

func NewCacheService() *CacheService {
	return &CacheService{
		redis: database.GetRedis(),
	}
}

const (
	PostCacheKeyPrefix = "post:"
	PostCacheTTL       = 5 * time.Minute
)

// GetPost retrieves a post from cache
func (cs *CacheService) GetPost(id uint) (*models.Post, error) {
	ctx := context.Background()
	key := fmt.Sprintf("%s%d", PostCacheKeyPrefix, id)

	val, err := cs.redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var post models.Post
	err = json.Unmarshal([]byte(val), &post)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

// SetPost stores a post in cache with TTL
func (cs *CacheService) SetPost(post *models.Post) error {
	ctx := context.Background()
	key := fmt.Sprintf("%s%d", PostCacheKeyPrefix, post.ID)

	postJSON, err := json.Marshal(post)
	if err != nil {
		return err
	}

	return cs.redis.Set(ctx, key, postJSON, PostCacheTTL).Err()
}

// InvalidatePost removes a post from cache
func (cs *CacheService) InvalidatePost(id uint) error {
	ctx := context.Background()
	key := fmt.Sprintf("%s%d", PostCacheKeyPrefix, id)

	return cs.redis.Del(ctx, key).Err()
}

// InvalidatePostsByPattern removes posts from cache by pattern
func (cs *CacheService) InvalidatePostsByPattern(pattern string) error {
	ctx := context.Background()

	keys, err := cs.redis.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return cs.redis.Del(ctx, keys...).Err()
	}

	return nil
}