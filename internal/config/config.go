package config

import (
	"fmt"
	"os"
)

type Config struct {
	Database      DatabaseConfig
	Redis         RedisConfig
	Elasticsearch ElasticsearchConfig
	Server        ServerConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type RedisConfig struct {
	Host string
	Port string
}

type ElasticsearchConfig struct {
	Host string
	Port string
}

type ServerConfig struct {
	Port string
}

func LoadConfig() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "blog_user"),
			Password: getEnv("DB_PASSWORD", "blog_password"),
			DBName:   getEnv("DB_NAME", "blog_db"),
		},
		Redis: RedisConfig{
			Host: getEnv("REDIS_HOST", "localhost"),
			Port: getEnv("REDIS_PORT", "6379"),
		},
		Elasticsearch: ElasticsearchConfig{
			Host: getEnv("ELASTICSEARCH_HOST", "localhost"),
			Port: getEnv("ELASTICSEARCH_PORT", "9200"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
	}
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Ho_Chi_Minh",
		d.Host, d.Port, d.User, d.Password, d.DBName)
}

func (r *RedisConfig) Address() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}

func (e *ElasticsearchConfig) URL() string {
	return fmt.Sprintf("http://%s:%s", e.Host, e.Port)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}