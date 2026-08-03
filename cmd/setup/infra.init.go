package setup

import (
	"cinema-ticketing-api/app/seed"
	"cinema-ticketing-api/config"
	"cinema-ticketing-api/pkg/database"
	"cinema-ticketing-api/pkg/mailer"
	"log"

	"gorm.io/gorm"
)

func InitInfra(cfg *config.Config) (*gorm.DB, mailer.Mailer) {

	db := database.ConnectDB(cfg.DB)
	if err := seed.SeedAdmin(db, cfg.Admin); err != nil {
		log.Fatalf("failed to seed admin: %v", err)
	}

	mailer := mailer.NewMailer(cfg.SMTP)

	return db, *mailer

}
