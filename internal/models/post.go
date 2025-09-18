package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Post struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"not null" binding:"required"`
	Content   string         `json:"content" gorm:"type:text;not null" binding:"required"`
	Tags      pq.StringArray `json:"tags" gorm:"type:text[]"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type ActivityLog struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	Action   string    `json:"action" gorm:"not null"`
	PostID   uint      `json:"post_id" gorm:"not null"`
	LoggedAt time.Time `json:"logged_at" gorm:"default:CURRENT_TIMESTAMP"`
}

type CreatePostRequest struct {
	Title   string   `json:"title" binding:"required"`
	Content string   `json:"content" binding:"required"`
	Tags    []string `json:"tags"`
}

type UpdatePostRequest struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type PostResponse struct {
	ID           uint      `json:"id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	Tags         []string  `json:"tags"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	RelatedPosts []Post    `json:"related_posts,omitempty"`
}

type SearchResponse struct {
	Posts []Post `json:"posts"`
	Total int64  `json:"total"`
}

// ElasticsearchPost represents post document in Elasticsearch
type ElasticsearchPost struct {
	ID      uint     `json:"id"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (p *Post) TableName() string {
	return "posts"
}

func (al *ActivityLog) TableName() string {
	return "activity_logs"
}