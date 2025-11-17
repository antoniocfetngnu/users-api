package database

import (
	"log"

	"github.com/antoniocfetngnu/users-api/config"
	"github.com/antoniocfetngnu/users-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) error {
	var err error

	// Configure GORM logger
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Connect to PostgreSQL
	DB, err = gorm.Open(postgres.Open(cfg.DatabaseURL), gormConfig)
	if err != nil {
		return err
	}

	log.Println("✅ Database connected successfully")

	// Auto-migrate models (creates tables if they don't exist)
	if err := DB.AutoMigrate(&models.User{}, &models.Follower{}); err != nil {
		return err
	}

	log.Println("✅ Database migrations completed")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
