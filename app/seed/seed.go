package seed

import (
	"cinema-ticketing-api/config"
	"cinema-ticketing-api/entities"
	"cinema-ticketing-api/pkg/password"
	"errors"
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
)

func SeedAdmin(db *gorm.DB, cfg config.AdminConfig) error {
	adminEmail := strings.TrimSpace(cfg.Email)
	adminPassword := cfg.Password

	if adminEmail == "" && adminPassword == "" {
		log.Println("Admin seeding skipped: ADMIN_EMAIL and ADMIN_PASSWORD are not set")
		return nil
	}

	if adminEmail == "" || adminPassword == "" {
		return errors.New("admin seeding requires both ADMIN_EMAIL and ADMIN_PASSWORD")
	}

	adminName := strings.TrimSpace(cfg.Name)
	if adminName == "" {
		adminName = "Administrator"
	}

	hashedPassword, err := password.HashPassword(adminPassword)
	if err != nil {
		return fmt.Errorf("hash admin password: %w", err)
	}

	var existingUser entities.User
	err = db.Where("email = ?", adminEmail).First(&existingUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			admin := entities.User{
				Name:     adminName,
				Email:    adminEmail,
				Password: hashedPassword,
				Role:     entities.RoleAdmin,
			}

			if createErr := db.Create(&admin).Error; createErr != nil {
				return fmt.Errorf("create admin user: %w", createErr)
			}

			log.Printf("Admin user seeded successfully for %s", adminEmail)
			return nil
		}

		return fmt.Errorf("find admin user: %w", err)
	}

	existingUser.Name = adminName
	existingUser.Email = adminEmail
	existingUser.Password = hashedPassword
	existingUser.Role = entities.RoleAdmin

	if err := db.Save(&existingUser).Error; err != nil {
		return fmt.Errorf("update admin user: %w", err)
	}

	log.Printf("Admin user ensured successfully for %s", adminEmail)
	return nil
}
