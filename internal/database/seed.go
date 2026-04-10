package database

import (
	"fmt"

	"avi_bd/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedDatabase seeds the database with initial data (optional)
func SeedDatabase(db *gorm.DB) error {
	// Check if users already exist
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		fmt.Println("Database already seeded")
		return nil
	}

	// Create test user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
	}

	if err := db.Create(&user).Error; err != nil {
		return err
	}

	// Create wallet for test user
	wallet := models.Wallet{
		UserID:  user.ID,
		Balance: 1000.0,
	}

	if err := db.Create(&wallet).Error; err != nil {
		return err
	}

	fmt.Println("Database seeded successfully")
	return nil
}
