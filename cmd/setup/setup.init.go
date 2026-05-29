package setup

import (
	"cinema-ticketing-api/internal/config"
	"cinema-ticketing-api/internal/middleware"
	"cinema-ticketing-api/internal/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func InitApp() {

	// Init konfigurasi
	cfg := config.Load()

	// Init db dan mailer
	db, mail := InitInfra(cfg)

	// Init modul controller dan router
	rc, txSvc, schSvc := InitModule(db, cfg, &mail)

	// Init scheduler
	InitScheduler(txSvc, schSvc)

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
