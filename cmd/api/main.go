package main

import (
	"cinema-ticketing-api/internal/config"
	"cinema-ticketing-api/internal/controller"
	"cinema-ticketing-api/internal/middleware"
	"cinema-ticketing-api/internal/repository"
	"cinema-ticketing-api/internal/routes"
	"cinema-ticketing-api/internal/service"
	"cinema-ticketing-api/pkg/database"
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

	// Inisialisasi Repositories
	userRepo := repository.NewUserRepository(db)
	studioRepo := repository.NewStudioRepository(db)
	movieRepo := repository.NewMovieRepository(db)
	seatRepo := repository.NewSeatRepository(db)
	scheduleRepo := repository.NewScheduleRepository(db)

	// Inisialisasi Services
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)
	studioService := service.NewStudioService(studioRepo)
	movieService := service.NewMovieService(movieRepo)
	seatService := service.NewSeatService(seatRepo)
	scheduleService := service.NewScheduleService(scheduleRepo)

	// Inisialisasi Controllers
	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService)
	studioController := controller.NewStudioController(studioService)
	movieController := controller.NewMovieController(movieService)
	seatController := controller.NewSeatController(seatService)
	scheduleController := controller.NewScheduleController(scheduleService)

	// Inisialisasi Router
	r := gin.Default()

	// Tambahan middleware aplikasi
	r.Use(middleware.CorsMiddleware())

	routes.SetupRoutes(r, &routes.RouteControllers{
		Auth:     authController,
		User:     userController,
		Studio:   studioController,
		Movie:    movieController,
		Seat:     seatController,
		Schedule: scheduleController,
	})

	// Init Port
	port := config.GetEnv("APP_PORT", ":8080")

	// Jalankan server
	r.Run(port)
}
