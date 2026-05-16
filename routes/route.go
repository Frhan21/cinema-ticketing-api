package routes

import (
	"cinema-ticketing-api/controller"
	"cinema-ticketing-api/middleware"
	"cinema-ticketing-api/repository"
	"cinema-ticketing-api/response"
	"cinema-ticketing-api/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {

	// User dan Auth
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)
	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService)

	// Studio
	studioRepo := repository.NewStudioRepository(db)
	studioService := service.NewStudioService(studioRepo)
	studioController := controller.NewStudioController(studioService)

	// Movie
	movieRepo := repository.NewMovieRepository(db)
	movieService := service.NewMovieService(movieRepo)
	movieController := controller.NewMovieController(movieService)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, response.SuccessResponse("Cinema Ticketing API is running", nil))
	})

	api := r.Group("/api/v1")

	// Authentication
	authRoute := api.Group("/auth")
	{
		authRoute.POST("/register", authController.Register)
		authRoute.POST("/login", authController.Login)
	}

	// Get data profile user dan Update Profile
	userRoute := api.Group("/user")
	userRoute.Use(middleware.AuthMiddleware())
	{
		userRoute.GET("/profile", userController.GetProfile)
		userRoute.PUT("/profile/", userController.UpdateProfile)
	}

	studioRoute := api.Group("/studio")
	studioRoute.Use(middleware.AuthMiddleware())
	{
		// Semua user boleh melihat list studio
		studioRoute.GET("/", studioController.GetAll)
		studioRoute.GET("/:id", studioController.GetByID)

		// Admin hanya boleh membuat, mengupdate, menghapus
		adminStudio := studioRoute.Group("/")
		adminStudio.Use(middleware.RoleMiddleware("admin"))
		{
			adminStudio.POST("/", studioController.Create)
			adminStudio.PUT("/:id", studioController.Update)
			adminStudio.DELETE("/:id", studioController.Delete)
		}
	}

	// Movie
	movieRoute := api.Group("/movie")
	{
		// Semua user boleh melihat list movie
		movieRoute.GET("/", movieController.GetAll)
		movieRoute.GET("/:id", movieController.GetByID)

		// Admin hanya boleh membuat, mengupdate, menghapus
		adminMovie := movieRoute.Group("/")
		adminMovie.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
		{
			adminMovie.POST("/", movieController.Create)
			adminMovie.PUT("/:id", movieController.Update)
			adminMovie.DELETE("/:id", movieController.Delete)
		}
	}
}
