package main

import (
	"blog-api/internal/config"
	"blog-api/internal/database"
	"blog-api/internal/handlers"
	"blog-api/internal/middleware"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize databases with retry logic
	if err := initializeDatabases(cfg); err != nil {
		log.Fatal("Failed to initialize databases:", err)
	}

	// Set up Gin router
	router := setupRouter()

	// Start server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func initializeDatabases(cfg *config.Config) error {
	// Initialize PostgreSQL with retry
	for i := 0; i < 5; i++ {
		if err := database.InitDB(&cfg.Database); err != nil {
			log.Printf("Failed to connect to PostgreSQL (attempt %d/5): %v", i+1, err)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	// Initialize Redis with retry
	for i := 0; i < 5; i++ {
		if err := database.InitRedis(&cfg.Redis); err != nil {
			log.Printf("Failed to connect to Redis (attempt %d/5): %v", i+1, err)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	// Initialize Elasticsearch with retry
	for i := 0; i < 5; i++ {
		if err := database.InitElasticsearch(&cfg.Elasticsearch); err != nil {
			log.Printf("Failed to connect to Elasticsearch (attempt %d/5): %v", i+1, err)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	log.Println("All databases initialized successfully")
	return nil
}

func setupRouter() *gin.Engine {
	// Set to release mode in production
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Add middleware
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())

	// Initialize handlers
	postHandler := handlers.NewPostHandler()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Blog API is running",
		})
	})

	// API routes
	v1 := router.Group("/api/v1")
	{
		posts := v1.Group("/posts")
		{
			posts.POST("", postHandler.CreatePost)
			posts.GET("", postHandler.GetAllPosts)
			posts.GET("/:id", postHandler.GetPost)
			posts.PUT("/:id", postHandler.UpdatePost)
			posts.DELETE("/:id", postHandler.DeletePost)
			posts.GET("/search-by-tag", postHandler.SearchPostsByTag)
			posts.GET("/search", postHandler.SearchPosts)
		}
	}

	return router
}