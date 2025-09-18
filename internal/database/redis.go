package database

import (
	"blog-api/internal/config"
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func ConnectRedis(cfg *config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Address(),
		Password: "", // no password
		DB:       0,  // default DB
	})

	// Test connection
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	log.Println("Connected to Redis successfully")
	return rdb, nil
}

func InitRedis(cfg *config.RedisConfig) error {
	var err error
	redisClient, err = ConnectRedis(cfg)
	return err
}

func GetRedis() *redis.Client {
	return redisClient
}