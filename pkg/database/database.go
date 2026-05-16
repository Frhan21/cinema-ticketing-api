package database

import (
	"cinema-ticketing-api/internal/config"
	"cinema-ticketing-api/internal/model"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	var errDb error
	var db *gorm.DB
	driver := config.GetEnv("DB_DRIVER", "mysql")

	switch driver {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.GetEnv("DB_USER", "root"), config.GetEnv("DB_PASSWORD", ""), config.GetEnv("DB_HOST", "127.0.0.1"), config.GetEnv("DB_PORT", "3306"), config.GetEnv("DB_NAME", "go_gin"))
		db, errDb = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", config.GetEnv("DB_HOST", "127.0.0.1"), config.GetEnv("DB_USER", "root"), config.GetEnv("DB_PASSWORD", ""), config.GetEnv("DB_NAME", "go_gin"), config.GetEnv("DB_PORT", "5432"))
		db, errDb = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	default:
		panic(fmt.Sprintf("Unsupported DB_DRIVER: %s", driver))
	}

	if errDb != nil {
		panic(fmt.Sprintf("Can't connect to database: %v", errDb))
	}

	sqlDB, errDb := db.DB()
	if errDb = sqlDB.Ping(); errDb != nil {
		panic(fmt.Sprintf("Database ping failed: %v", errDb))
	}

	log.Println("Successfully connected to the database")
	return db
}

func Migration(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.User{},
		&model.Studio{},
		&model.Movie{},
	)

	if err != nil {
		log.Println("Error Migration : ", err)
	}
}
