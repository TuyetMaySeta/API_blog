package database

import (
	"blog-api/internal/config"
	"blog-api/internal/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectPostgres(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	log.Println("Connected to PostgreSQL successfully")

	// Auto migrate tables
	err = db.AutoMigrate(&models.Post{}, &models.ActivityLog{})
	if err != nil {
		log.Printf("Error migrating tables: %v", err)
		return nil, err
	}

	log.Println("Tables migrated successfully")
	return db, nil
}

func GetDB() *gorm.DB {
	return db
}

var db *gorm.DB

func InitDB(cfg *config.DatabaseConfig) error {
	var err error
	db, err = ConnectPostgres(cfg)
	return err
}