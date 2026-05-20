// Package main adalah entry point dari Cinema Ticketing API.
//
// @title           Cinema Ticketing API
// @version         1.0
// @description     REST API untuk sistem manajemen tiket bioskop. Mendukung manajemen studio, film, jadwal tayang, kursi, booking tiket, transaksi, promo diskon, dan laporan penjualan.
// @termsOfService  http://swagger.io/terms/
//
// @contact.name    Cinema Ticketing Support
// @contact.email   support@cinemticket.com
//
// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT
//
// @host      localhost:8080
// @BasePath  /api/v1
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Masukkan token dengan format: Bearer {token}
package main

import (
	"cinema-ticketing-api/internal/config"
	"cinema-ticketing-api/internal/controller"
	"cinema-ticketing-api/internal/middleware"
	"cinema-ticketing-api/internal/repository"
	"cinema-ticketing-api/internal/routes"
	"cinema-ticketing-api/internal/service"
	"cinema-ticketing-api/pkg/database"
	"cinema-ticketing-api/pkg/mailer"
	"cinema-ticketing-api/pkg/scheduler"
	"log"

	_ "cinema-ticketing-api/docs"

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

	// Inisialisasi Mailer
	mail := mailer.NewMailer(cfg.SMTP)

	// Inisialisasi Repositories
	userRepo := repository.NewUserRepository(db)
	studioRepo := repository.NewStudioRepository(db)
	movieRepo := repository.NewMovieRepository(db)
	seatRepo := repository.NewSeatRepository(db)
	scheduleRepo := repository.NewScheduleRepository(db)
	ticketRepo := repository.NewTicketRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	transactionItemRepo := repository.NewTransactionItemRepository(db)
	promoRepo := repository.NewPromoRepository(db)
	reportRepo := repository.NewReportRepository(db)

	// Inisialisasi Services
	authService := service.NewAuthService(userRepo, cfg.JWT)
	userService := service.NewUserService(userRepo)
	studioService := service.NewStudioService(studioRepo)
	movieService := service.NewMovieService(movieRepo)
	seatService := service.NewSeatService(seatRepo)
	scheduleService := service.NewScheduleService(scheduleRepo)
	ticketService := service.NewTicketService(ticketRepo, seatRepo, transactionRepo, transactionItemRepo, scheduleRepo, promoRepo)
	transactionService := service.NewTransactionService(transactionRepo, ticketRepo, mail)
	promoService := service.NewPromoService(promoRepo)
	reportService := service.NewReportService(reportRepo)

	// Inisialisasi Controllers
	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService)
	studioController := controller.NewStudioController(studioService)
	movieController := controller.NewMovieController(movieService)
	seatController := controller.NewSeatController(seatService)
	scheduleController := controller.NewScheduleController(scheduleService)
	ticketController := controller.NewTicketController(ticketService)
	transactionController := controller.NewTransactionController(transactionService)
	promoController := controller.NewPromoController(promoService)
	reportController := controller.NewReportController(reportService)

	// Start background scheduler (auto-cancel & film reminder)
	sched := scheduler.NewScheduler(db, mail)
	sched.Start()

	// Inisialisasi Router
	r := gin.Default()

	// Tambahan middleware aplikasi
	r.Use(middleware.CorsMiddleware())

	routes.SetupRoutes(r, &routes.RouteControllers{
		Auth:        authController,
		User:        userController,
		Studio:      studioController,
		Movie:       movieController,
		Seat:        seatController,
		Schedule:    scheduleController,
		Ticket:      ticketController,
		Transaction: transactionController,
		Promo:       promoController,
		Report:      reportController,
	}, cfg.JWT)

	// Jalankan server
	log.Printf("Server running on port %s", cfg.App.Port)
	log.Printf("Swagger UI: http://localhost%s/swagger/index.html", cfg.App.Port)
	r.Run(cfg.App.Port)
}
