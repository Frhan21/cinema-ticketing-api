package main

import (
	"cinema-ticketing-api/internal/config"
	"cinema-ticketing-api/pkg/database"
	"cinema-ticketing-api/internal/middleware"
	"cinema-ticketing-api/internal/routes"
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
	r := gin.Default()

	// Tambahan middleware aplikasi
	r.Use(middleware.CorsMiddleware())

	routes.SetupRoutes(r, db)

	// Init Port
	port := config.GetEnv("APP_PORT", ":8080")

	// Jalankan server
	r.Run(port)
}
