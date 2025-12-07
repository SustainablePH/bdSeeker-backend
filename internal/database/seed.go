package database

import (
	"log"

	"github.com/bishworup11/bdSeeker-backend/internal/models"
	"github.com/bishworup11/bdSeeker-backend/pkg/utils"
)

// SeedAdminUser creates a default admin user if it doesn't exist
func SeedAdminUser() error {
	var count int64
	DB.Model(&models.User{}).Where("role = ?", "admin").Count(&count)

	if count > 0 {
		log.Println("✓ Admin user already exists")
		return nil
	}

	// Create default admin user
	hashedPassword, err := utils.HashPassword("admin123")
	if err != nil {
		return err
	}

	admin := &models.User{
		Email:        "admin@bdseeker.com",
		PasswordHash: hashedPassword,
		FullName:     "System Administrator",
		Role:         "admin",
	}

	if err := DB.Create(admin).Error; err != nil {
		return err
	}

	log.Println("✓ Default admin user created successfully")
	log.Println("  Email: admin@bdseeker.com")
	log.Println("  Password: admin123")
	log.Println("  ⚠️  IMPORTANT: Change the admin password in production!")

	return nil
}
