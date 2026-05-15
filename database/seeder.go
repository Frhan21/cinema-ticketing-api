package database

import (
	"cinema-ticketing-api/config"
	"cinema-ticketing-api/models"
	"cinema-ticketing-api/utils"
	"errors"
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
)

func SeedAdmin(db *gorm.DB) error {
	adminEmail := strings.TrimSpace(config.GetEnv("ADMIN_EMAIL", ""))
	adminPassword := config.GetEnv("ADMIN_PASSWORD", "")

	if adminEmail == "" && adminPassword == "" {
		log.Println("Admin seeding skipped: ADMIN_EMAIL and ADMIN_PASSWORD are not set")
		return nil
	}

	if adminEmail == "" || adminPassword == "" {
		return errors.New("admin seeding requires both ADMIN_EMAIL and ADMIN_PASSWORD")
	}

	adminName := strings.TrimSpace(config.GetEnv("ADMIN_NAME", "Administrator"))
	if adminName == "" {
		adminName = "Administrator"
	}

	hashedPassword, err := utils.HashPassword(adminPassword)
	if err != nil {
		return fmt.Errorf("hash admin password: %w", err)
	}

	var user models.User
	err = db.Where("email = ?", adminEmail).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			admin := models.User{
				Name:     adminName,
				Email:    adminEmail,
				Password: hashedPassword,
				Role:     models.AdminRole,
			}

			if createErr := db.Create(&admin).Error; createErr != nil {
				return fmt.Errorf("create admin user: %w", createErr)
			}

			log.Printf("Admin user seeded successfully for %s", adminEmail)
			return nil
		}

		return fmt.Errorf("find admin user: %w", err)
	}

	user.Name = adminName
	user.Email = adminEmail
	user.Password = hashedPassword
	user.Role = models.AdminRole

	if err := db.Save(&user).Error; err != nil {
		return fmt.Errorf("update admin user: %w", err)
	}

	log.Printf("Admin user ensured successfully for %s", adminEmail)
	return nil
}
