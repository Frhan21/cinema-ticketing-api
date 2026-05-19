package database

import (
	"cinema-ticketing-api/internal/config"
	"cinema-ticketing-api/internal/model"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(cfg config.DBConfig) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Asia/Jakarta",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Can't connect to database: %v", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("Failed to get sql.DB: %v", err))
	}
	if err = sqlDB.Ping(); err != nil {
		panic(fmt.Sprintf("Database ping failed: %v", err))
	}

	runMigration(cfg)

	log.Println("Successfully connected to the database")
	return db
}

func runMigration(cfg config.DBConfig) {
	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=require",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
	)

	m, err := migrate.New("file://./database/migrations", databaseURL)
	if err != nil {
		log.Fatal("Migration initialization failed:", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Database migration executed successfully")
}

func AutoMigration(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.User{},
		&model.Studio{},
		&model.Movie{},
		&model.Seat{},
		&model.Schedule{},
		&model.Ticket{},
		&model.Transaction{},
		&model.TransactionItem{},
	)

	if err != nil {
		log.Println("Error AutoMigration:", err)
	}
}
