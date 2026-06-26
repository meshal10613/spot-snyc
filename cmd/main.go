package main

import (
	"fmt"
	"log"
	"net/http"
	"spot-sync/config"
	"spot-sync/handler"
	"spot-sync/models"
	"spot-sync/repository"
	"spot-sync/routes"
	"spot-sync/service"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

// CustomValidator implements echo.Validator using go-playground/validator.
type CustomValidator struct {
	validator *validator.Validate
}

// Validate validates structs using go-playground/validator tags.
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	// ── Load Configuration ────────────────────────────────────────────────
	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatal("Failed to load configuration: ", err)
	}

	// ── Connect Database ──────────────────────────────────────────────────
	db := config.ConnectDatabase(cfg)

	// ── Run Migrations ────────────────────────────────────────────────────
	fmt.Println("⏳ Running database migrations...")
	if err := models.Migrate(db); err != nil {
		log.Fatal("❌ Migration failed: ", err)
	}
	fmt.Println("✅ Migrations completed")

	// ── Initialize Echo ───────────────────────────────────────────────────
	e := echo.New()

	// Validator — integrated with Echo so c.Validate() works globally
	e.Validator = &CustomValidator{validator: validator.New()}

	// Global middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodPatch,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
	}))

	// ── Health Check ──────────────────────────────────────────────────────
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "SpotSync API is running 🚗",
		})
	})

	// ── Dependency Injection (Manual Wiring) ──────────────────────────────
	// Repositories
	authRepo := repository.NewAuthRepository(db)
	zoneRepo := repository.NewZoneRepository(db)
	reservationRepo := repository.NewReservationRepository(db)

	// Services
	authSvc := service.NewAuthService(authRepo, cfg)
	zoneSvc := service.NewZoneService(zoneRepo)
	reservationSvc := service.NewReservationService(reservationRepo, zoneRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authSvc)
	zoneHandler := handler.NewZoneHandler(zoneSvc)
	reservationHandler := handler.NewReservationHandler(reservationSvc)

	// ── Register Routes ───────────────────────────────────────────────────
	routes.RegisterRoutes(e, authHandler, zoneHandler, reservationHandler, cfg.JWTSecret)

	// ── Start Server ──────────────────────────────────────────────────────
	fmt.Printf("🚀 Server running at http://localhost:%s\n", cfg.Port)
	if err := e.Start(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
		log.Fatal("❌ Failed to start server: ", err)
	}
}
