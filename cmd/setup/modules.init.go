package setup

import (
	"cinema-ticketing-api/internal/config"
	"cinema-ticketing-api/internal/controller"
	"cinema-ticketing-api/internal/repository"
	"cinema-ticketing-api/internal/routes"
	"cinema-ticketing-api/internal/service"
	"cinema-ticketing-api/pkg/mailer"

	"gorm.io/gorm"
)

func InitModule(db *gorm.DB, cfg *config.Config, mail *mailer.Mailer) (*routes.RouteControllers, service.TransactionService, service.ScheduleService) {

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
	scheduleService := service.NewScheduleService(scheduleRepo, ticketRepo, mail)
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

	return &routes.RouteControllers{
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
	}, transactionService, scheduleService
}
