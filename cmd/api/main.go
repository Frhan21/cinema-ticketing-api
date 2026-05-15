package main

import (
	"cinema-ticketing-api/config"
	"cinema-ticketing-api/database"
	"cinema-ticketing-api/middleware"
	"cinema-ticketing-api/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load Konfigurasi dari .env
	config.LoadEnv()

	// Konek ke DB
	db := database.ConnectDB()
	// Auto migrasi db
	database.Migration(db)
	if err := database.SeedAdmin(db); err != nil {
		log.Fatalf("failed to seed admin: %v", err)
	}

	// Inisialisasi Router
	r := gin.New()

	// Gunakan standard middleware: Panic Recovery, Custom Logger, & CORS
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.CorsMiddleware())

	routes.SetupRoutes(r, db)

	// Init Port
	port := config.GetEnv("APP_PORT", ":8080")

	// Jalankan server
	r.Run(port)
}
