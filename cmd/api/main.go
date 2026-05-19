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
	cfg := config.Load()

	// Konek ke DB (termasuk auto-migration)
	db := database.ConnectDB(cfg.DB)

	// Seed admin dari konfigurasi
	if err := database.SeedAdmin(db, cfg.Admin); err != nil {
		log.Fatalf("failed to seed admin: %v", err)
	}

	// Inisialisasi Repositories
	userRepo := repository.NewUserRepository(db)
	studioRepo := repository.NewStudioRepository(db)
	movieRepo := repository.NewMovieRepository(db)
	seatRepo := repository.NewSeatRepository(db)
	scheduleRepo := repository.NewScheduleRepository(db)

	// Inisialisasi Services
	authService := service.NewAuthService(userRepo, cfg.JWT)
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
	}, cfg.JWT)

	// Jalankan server
	log.Printf("Server running on port %s", cfg.App.Port)
	r.Run(cfg.App.Port)
}
