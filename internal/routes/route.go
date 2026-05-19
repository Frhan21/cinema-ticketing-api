package routes

import (
	"cinema-ticketing-api/internal/controller"
	"cinema-ticketing-api/internal/middleware"
	"cinema-ticketing-api/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RouteControllers struct {
	Auth     *controller.AuthController
	User     *controller.UserController
	Studio   controller.StudioController
	Movie    controller.MovieController
	Seat     controller.SeatController
	Schedule controller.ScheduleController
}

func SetupRoutes(r *gin.Engine, ctrl *RouteControllers) {

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, response.SuccessResponse("Cinema Ticketing API is running", nil))
	})

	api := r.Group("/api/v1")

	// Test Send Email
	r.POST("/api/v1/test-email", controller.TestSendEmail)

	// Authentication
	authRoute := api.Group("/auth")
	{
		authRoute.POST("/register", ctrl.Auth.Register)
		authRoute.POST("/login", ctrl.Auth.Login)
	}

	// Get data profile user dan Update Profile
	userRoute := api.Group("/user")
	userRoute.Use(middleware.AuthMiddleware())
	{
		userRoute.GET("/profile", ctrl.User.GetProfile)
		userRoute.PUT("/profile/", ctrl.User.UpdateProfile)
	}

	// Studio
	studioRoute := api.Group("/studio")
	studioRoute.Use(middleware.AuthMiddleware())
	{
		// Semua user boleh melihat list studio
		studioRoute.GET("/", ctrl.Studio.GetAll)
		studioRoute.GET("/:id", ctrl.Studio.GetByID)

		// Admin hanya boleh membuat, mengupdate, menghapus
		adminStudio := studioRoute.Group("/")
		adminStudio.Use(middleware.RoleMiddleware("admin"))
		{
			adminStudio.POST("/", ctrl.Studio.Create)
			adminStudio.PUT("/:id", ctrl.Studio.Update)
			adminStudio.DELETE("/:id", ctrl.Studio.Delete)
		}
	}

	// Movie
	movieRoute := api.Group("/movie")
	{
		// Semua user boleh melihat list movie
		movieRoute.GET("/", ctrl.Movie.GetAll)
		movieRoute.GET("/:id", ctrl.Movie.GetByID)

		// Admin hanya boleh membuat, mengupdate, menghapus
		adminMovie := movieRoute.Group("/")
		adminMovie.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
		{
			adminMovie.POST("/", ctrl.Movie.Create)
			adminMovie.PUT("/:id", ctrl.Movie.Update)
			adminMovie.DELETE("/:id", ctrl.Movie.Delete)
		}
	}

	// Seat
	seatRoute := api.Group("/seat")
	seatRoute.Use(middleware.AuthMiddleware())
	{
		// Semua user boleh melihat list seat
		seatRoute.GET("/", ctrl.Seat.FindAll)
		seatRoute.GET("/:id", ctrl.Seat.FindByID)
		seatRoute.GET("/studio/:studio_id", ctrl.Seat.FindByStudioID)

		// Admin hanya boleh membuat, mengupdate, menghapus
		adminSeat := seatRoute.Group("/")
		adminSeat.Use(middleware.RoleMiddleware("admin"))
		{
			adminSeat.POST("/", ctrl.Seat.Create)
			adminSeat.PUT("/:id", ctrl.Seat.Update)
			adminSeat.DELETE("/:id", ctrl.Seat.Delete)
		}
	}

	// Schedule
	scheduleRoute := api.Group("/schedule")
	scheduleRoute.Use(middleware.AuthMiddleware())
	{
		// Semua user boleh melihat list schedule
		scheduleRoute.GET("/", ctrl.Schedule.FindAll)
		scheduleRoute.GET("/:id", ctrl.Schedule.FindByID)

		// Admin hanya boleh membuat, mengupdate, menghapus
		adminSchedule := scheduleRoute.Group("/")
		adminSchedule.Use(middleware.RoleMiddleware("admin"))
		{
			adminSchedule.POST("/", ctrl.Schedule.Create)
			adminSchedule.PUT("/:id", ctrl.Schedule.Update)
			adminSchedule.DELETE("/:id", ctrl.Schedule.Delete)
		}
	}
}
