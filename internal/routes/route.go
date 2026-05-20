package routes

import (
	"cinema-ticketing-api/internal/config"
	"cinema-ticketing-api/internal/controller"
	"cinema-ticketing-api/internal/middleware"
	"cinema-ticketing-api/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RouteControllers struct {
	Auth        *controller.AuthController
	User        *controller.UserController
	Studio      controller.StudioController
	Movie       controller.MovieController
	Seat        controller.SeatController
	Schedule    controller.ScheduleController
	Ticket      controller.TicketController
	Transaction controller.TransactionController
	Promo       controller.PromoController
	Report      controller.ReportController
}

func SetupRoutes(r *gin.Engine, ctrl *RouteControllers, jwtCfg config.JWTConfig) {

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, response.SuccessResponse("Cinema Ticketing API is running", nil))
	})

	api := r.Group("/api/v1")

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

	// Schedule — authenticated (read) | admin (write)
	scheduleRoute := api.Group("/schedule")
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
		transactionRoute.POST("/:transaction_id/pay", ctrl.Transaction.PayTransaction)
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
