package setup

import (
	"cinema-ticketing-api/config"
	"cinema-ticketing-api/interface/http/middleware"
	"cinema-ticketing-api/interface/http/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func InitApp() {

	// Init konfigurasi
	cfg := config.Load()

	// Init db dan mailer
	db, mail := InitInfra(cfg)

	// Init modul controller dan router
	rc, scheduler, err := InitModule(db, cfg, &mail)
	if err != nil {
		log.Fatalf("Failed to initialize modules: %v", err)
	}

	// Init scheduler
	scheduler.Start()

	// init gin Default
	r := gin.Default()

	// Tambahan middleware aplikasi
	r.Use(middleware.ErrorHandlerMiddleware())
	r.Use(middleware.CorsMiddleware())

	routes.SetupRoutes(r, rc, cfg.JWT)

	// Jalankan server
	log.Printf("Server running on port %s", cfg.App.Port)
	log.Printf("Swagger UI: http://localhost%s/swagger/index.html", cfg.App.Port)
	r.Run(cfg.App.Port)

}
