package routes

import (
	"cinema-ticketing-api/config"
	"cinema-ticketing-api/interface/http/handler"
	"cinema-ticketing-api/interface/http/middleware"
	"cinema-ticketing-api/response"
	"net/http"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

type RouteControllers struct {
	Auth        *handler.AuthController
	User        *handler.UserController
	Studio      handler.StudioController
	Movie       handler.MovieController
	Seat        handler.SeatController
	Schedule    handler.ScheduleController
	Ticket      handler.TicketController
	Transaction handler.TransactionController
	Payment     handler.PaymentController
	Promo       handler.PromoController
	Report      handler.ReportController
}

func SetupRoutes(r *gin.Engine, ctrl *RouteControllers, jwtCfg config.JWTConfig) {

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, response.SuccessResponse("Cinema Ticketing API is running", nil))
	})

	// Swagger UI — accessible at /swagger/index.html
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")

	// Payment notification is public because it is called by Midtrans.
	paymentRoute := api.Group("/payment")
	{
		paymentRoute.POST("/notification", ctrl.Payment.HandleNotification)
	}

	// Authentication — public
	authRoute := api.Group("/auth")
	{
		authRoute.POST("/register", ctrl.Auth.Register)
		authRoute.POST("/login", ctrl.Auth.Login)
	}

	// User — authenticated
	userRoute := api.Group("/user")
	userRoute.Use(middleware.AuthMiddleware(jwtCfg))
	{
		userRoute.GET("/profile", ctrl.User.GetProfile)
		userRoute.PUT("/profile", ctrl.User.UpdateProfile)
	}

	// Studio — public (read) | admin (write)
	studioRoute := api.Group("/studio")
	studioRoute.Use(middleware.AuthMiddleware(jwtCfg))
	{
		studioRoute.GET("/", ctrl.Studio.GetAll)
		studioRoute.GET("/:id", ctrl.Studio.GetByID)

		adminStudio := studioRoute.Group("/")
		adminStudio.Use(middleware.RoleMiddleware("admin"))
		{
			adminStudio.POST("/", ctrl.Studio.Create)
			adminStudio.PUT("/:id", ctrl.Studio.Update)
			adminStudio.DELETE("/:id", ctrl.Studio.Delete)
		}
	}

	// Movie — public (read) | admin (write)
	movieRoute := api.Group("/movie")
	{
		movieRoute.GET("/", ctrl.Movie.GetAll)
		movieRoute.GET("/:id", ctrl.Movie.GetByID)

		adminMovie := movieRoute.Group("/")
		adminMovie.Use(middleware.AuthMiddleware(jwtCfg), middleware.RoleMiddleware("admin"))
		{
			adminMovie.POST("/import", ctrl.Movie.Import)
			adminMovie.POST("/", ctrl.Movie.Create)
			adminMovie.PUT("/:id", ctrl.Movie.Update)
			adminMovie.DELETE("/:id", ctrl.Movie.Delete)
		}
	}

	// Seat — authenticated (read) | admin (write)
	seatRoute := api.Group("/seat")
	seatRoute.Use(middleware.AuthMiddleware(jwtCfg))
	{
		seatRoute.GET("/", ctrl.Seat.FindAll)
		seatRoute.GET("/:id", ctrl.Seat.FindByID)
		seatRoute.GET("/studio/:studio_id", ctrl.Seat.FindByStudioID)

		adminSeat := seatRoute.Group("/")
		adminSeat.Use(middleware.RoleMiddleware("admin"))
		{
			adminSeat.POST("/", ctrl.Seat.Create)
			adminSeat.PUT("/:id", ctrl.Seat.Update)
			adminSeat.DELETE("/:id", ctrl.Seat.Delete)
		}
	}

	// Schedule — public upcoming reads | authenticated detail reads | admin writes
	scheduleRoute := api.Group("/schedule")
	scheduleRoute.GET("/upcoming", ctrl.Schedule.FindUpcoming)
	scheduleRoute.GET("/movie/tmdb/:tmdb_id", ctrl.Schedule.FindUpcomingByTMDBID)
	scheduleRoute.Use(middleware.AuthMiddleware(jwtCfg))
	{
		scheduleRoute.GET("/", ctrl.Schedule.FindAll)
		scheduleRoute.GET("/:id", ctrl.Schedule.FindByID)

		adminSchedule := scheduleRoute.Group("/")
		adminSchedule.Use(middleware.RoleMiddleware("admin"))
		{
			adminSchedule.POST("/", ctrl.Schedule.Create)
			adminSchedule.PUT("/:id", ctrl.Schedule.Update)
			adminSchedule.DELETE("/:id", ctrl.Schedule.Delete)
		}
	}

	// Ticket — authenticated
	ticketRoute := api.Group("/ticket")
	ticketRoute.Use(middleware.AuthMiddleware(jwtCfg))
	{
		ticketRoute.POST("/", ctrl.Ticket.BookTicket)
		ticketRoute.GET("/history", ctrl.Ticket.GetUserHistory)
		ticketRoute.GET("/available-seats/:schedule_id", ctrl.Ticket.GetAvailableSeats)
	}

	// Transaction — authenticated (user) | admin (read all)
	transactionRoute := api.Group("/transaction")
	transactionRoute.Use(middleware.AuthMiddleware(jwtCfg))
	{
		transactionRoute.GET("/:transaction_id", ctrl.Transaction.GetByID)
		transactionRoute.POST("/:transaction_id/pay", ctrl.Payment.CreatePayment)
		transactionRoute.POST("/:transaction_id/cancel", ctrl.Transaction.CancelTransaction)

		adminTransaction := transactionRoute.Group("/")
		adminTransaction.Use(middleware.RoleMiddleware("admin"))
		{
			adminTransaction.GET("/", ctrl.Transaction.GetAll)
		}
	}

	// Promo — public (validate) | admin (CRUD)
	promoRoute := api.Group("/promo")
	promoRoute.Use(middleware.AuthMiddleware(jwtCfg))
	{
		// User dapat memvalidasi kode promo
		promoRoute.GET("/validate/:code", ctrl.Promo.ValidateCode)

		adminPromo := promoRoute.Group("/")
		adminPromo.Use(middleware.RoleMiddleware("admin"))
		{
			adminPromo.POST("/", ctrl.Promo.Create)
			adminPromo.GET("/", ctrl.Promo.GetAll)
			adminPromo.GET("/:id", ctrl.Promo.GetByID)
			adminPromo.PUT("/:id", ctrl.Promo.Update)
			adminPromo.DELETE("/:id", ctrl.Promo.Delete)
		}
	}

	// Report — admin only
	reportRoute := api.Group("/report")
	reportRoute.Use(middleware.AuthMiddleware(jwtCfg), middleware.RoleMiddleware("admin"))
	{
		reportRoute.GET("/daily", ctrl.Report.GetDailyReport)
		reportRoute.GET("/monthly", ctrl.Report.GetMonthlyReport)
	}
}
