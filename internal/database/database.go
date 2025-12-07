package database

import (
	"fmt"
	"log"

	"github.com/bishworup11/bdSeeker-backend/internal/config"
	"github.com/bishworup11/bdSeeker-backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect establishes a connection to the PostgreSQL database
func Connect(cfg *config.Config) error {
	dsn := cfg.GetDSN()

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	log.Println("✓ Database connection established successfully")
	return nil
}

// Migrate runs auto-migration for all models
func Migrate() error {
	log.Println("Running database migrations...")

	// Migrate in order of dependencies
	err := DB.AutoMigrate(
		// Base models
		&models.User{},
		&models.Technology{},
		&models.ProgrammingLanguage{},

		// Company models
		&models.CompanyProfile{},
		&models.CompanyRating{},
		&models.CompanyReview{},
		&models.CompanyReviewReaction{},
		&models.CompanyReviewComment{},
		&models.CompanyReviewReply{},

		// Developer models
		&models.DeveloperProfile{},
		&models.DeveloperExperience{},
		&models.DeveloperEducation{},
		&models.DeveloperCertificate{},

		// Job models
		&models.JobPost{},
		&models.PostReaction{},
		&models.PostComment{},
		&models.CommentReply{},

		// Report model
		&models.UserReport{},
	)

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("✓ Database migrations completed successfully")
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
