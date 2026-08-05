package setup

import (
	authservice "cinema-ticketing-api/app/auth/service"
	movierepository "cinema-ticketing-api/app/movie/repository"
	movieservice "cinema-ticketing-api/app/movie/service"
	paymentgateway "cinema-ticketing-api/app/payment/gateway"
	paymentrepository "cinema-ticketing-api/app/payment/repository"
	paymentservice "cinema-ticketing-api/app/payment/service"
	promorepository "cinema-ticketing-api/app/promo/repository"
	promoservice "cinema-ticketing-api/app/promo/service"
	reportrepository "cinema-ticketing-api/app/report/repository"
	reportservice "cinema-ticketing-api/app/report/service"
	schedulerepository "cinema-ticketing-api/app/schedule/repository"
	scheduleservice "cinema-ticketing-api/app/schedule/service"
	seatrepository "cinema-ticketing-api/app/seat/repository"
	seatservice "cinema-ticketing-api/app/seat/service"
	studiorepository "cinema-ticketing-api/app/studio/repository"
	studioservice "cinema-ticketing-api/app/studio/service"
	ticketrepository "cinema-ticketing-api/app/ticket/repository"
	ticketservice "cinema-ticketing-api/app/ticket/service"
	userrepository "cinema-ticketing-api/app/user/repository"
	userservice "cinema-ticketing-api/app/user/service"
	"cinema-ticketing-api/config"
	"cinema-ticketing-api/interface/http/handler"
	"cinema-ticketing-api/interface/http/routes"
	"cinema-ticketing-api/job"
	"cinema-ticketing-api/pkg/mailer"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func InitModule(db *gorm.DB, cfg *config.Config, mail *mailer.Mailer) (*routes.RouteControllers, *job.Scheduler, error) {

	// Inisialisasi Repositories
	userRepo := userrepository.NewUserRepository(db)
	studioRepo := studiorepository.NewStudioRepository(db)
	movieRepo := movierepository.NewMovieRepository(db)
	seatRepo := seatrepository.NewSeatRepository(db)
	scheduleRepo := schedulerepository.NewScheduleRepository(db)
	ticketRepo := ticketrepository.NewTicketRepository(db)
	transactionRepo := ticketrepository.NewTransactionRepository(db)
	promoRepo := promorepository.NewPromoRepository(db)
	reportRepo := reportrepository.NewReportRepository(db)
	paymentRepo := paymentrepository.NewPaymentRepository(db)

	paymentGateway, err := paymentgateway.NewMidtransGateway(cfg.Midtrans)
	if err != nil {
		return nil, nil, fmt.Errorf("initialize Midtrans gateway: %w", err)
	}

	// Inisialisasi Services
	authService := authservice.NewAuthService(userRepo, cfg.JWT)
	userService := userservice.NewUserService(userRepo)
	studioService := studioservice.NewStudioService(studioRepo)
	movieService := movieservice.NewMovieService(movieRepo)
	seatService := seatservice.NewSeatService(seatRepo)
	scheduleService := scheduleservice.NewScheduleService(scheduleRepo)
	ticketService := ticketservice.NewTicketService(ticketRepo, seatRepo, transactionRepo, scheduleRepo, promoRepo)
	paymentExpiry := time.Duration(cfg.Midtrans.ExpiryMinutes) * time.Minute
	transactionService := ticketservice.NewTransactionService(transactionRepo, ticketRepo, mail, paymentExpiry)
	promoService := promoservice.NewPromoService(promoRepo)
	reportService := reportservice.NewReportService(reportRepo)
	paymentService := paymentservice.NewPaymentService(paymentRepo, transactionRepo, paymentGateway, mail)

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
	paymentController := handler.NewPaymentController(paymentService)

	// Inisialisasi Scheduler
	scheduler := job.NewScheduler(transactionService, scheduleRepo, ticketRepo, mail, paymentExpiry)

	return &routes.RouteControllers{
		Auth:        authController,
		User:        userController,
		Studio:      studioController,
		Movie:       movieController,
		Seat:        seatController,
		Schedule:    scheduleController,
		Ticket:      ticketController,
		Transaction: transactionController,
		Payment:     paymentController,
		Promo:       promoController,
		Report:      reportController,
	}, scheduler, nil
}
