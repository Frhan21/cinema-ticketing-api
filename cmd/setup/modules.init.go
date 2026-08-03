package setup

import (
	"cinema-ticketing-api/app/auth"
	"cinema-ticketing-api/app/movie"
	"cinema-ticketing-api/app/promo"
	"cinema-ticketing-api/app/report"
	"cinema-ticketing-api/app/schedule"
	"cinema-ticketing-api/app/seat"
	"cinema-ticketing-api/app/studio"
	"cinema-ticketing-api/app/ticket"
	"cinema-ticketing-api/app/user"
	"cinema-ticketing-api/config"
	"cinema-ticketing-api/interface/http/handler"
	"cinema-ticketing-api/interface/http/routes"
	"cinema-ticketing-api/job"
	"cinema-ticketing-api/pkg/mailer"

	"gorm.io/gorm"
)

func InitModule(db *gorm.DB, cfg *config.Config, mail *mailer.Mailer) (*routes.RouteControllers, *job.Scheduler) {

	// Inisialisasi Repositories
	userRepo := user.NewUserRepository(db)
	studioRepo := studio.NewStudioRepository(db)
	movieRepo := movie.NewMovieRepository(db)
	seatRepo := seat.NewSeatRepository(db)
	scheduleRepo := schedule.NewScheduleRepository(db)
	ticketRepo := ticket.NewTicketRepository(db)
	transactionRepo := ticket.NewTransactionRepository(db)
	transactionItemRepo := ticket.NewTransactionItemRepository(db)
	promoRepo := promo.NewPromoRepository(db)
	reportRepo := report.NewReportRepository(db)

	// Inisialisasi Services
	authService := auth.NewAuthService(userRepo, cfg.JWT)
	userService := user.NewUserService(userRepo)
	studioService := studio.NewStudioService(studioRepo)
	movieService := movie.NewMovieService(movieRepo)
	seatService := seat.NewSeatService(seatRepo)
	scheduleService := schedule.NewScheduleService(scheduleRepo)
	ticketService := ticket.NewTicketService(ticketRepo, seatRepo, transactionRepo, transactionItemRepo, scheduleRepo, promoRepo)
	transactionService := ticket.NewTransactionService(transactionRepo, ticketRepo, mail)
	promoService := promo.NewPromoService(promoRepo)
	reportService := report.NewReportService(reportRepo)

	// Inisialisasi Handlers (controllers)
	authController := handler.NewAuthController(authService)
	userController := handler.NewUserController(userService)
	studioController := handler.NewStudioController(studioService)
	movieController := handler.NewMovieController(movieService)
	seatController := handler.NewSeatController(seatService)
	scheduleController := handler.NewScheduleController(scheduleService)
	ticketController := handler.NewTicketController(ticketService)
	transactionController := handler.NewTransactionController(transactionService)
	promoController := handler.NewPromoController(promoService)
	reportController := handler.NewReportController(reportService)

	// Inisialisasi Scheduler
	scheduler := job.NewScheduler(transactionService, scheduleRepo, ticketRepo, mail)

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
	}, scheduler
}
