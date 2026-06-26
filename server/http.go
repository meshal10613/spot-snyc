package server

import (
	"net/http"
	"spot-sync/config"
	"spot-sync/handler"
	"spot-sync/httpresponse"
	"spot-sync/repository"
	"spot-sync/routes"
	"spot-sync/service"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

// customValidator implements echo.Validator using go-playground/validator.
type customValidator struct {
	validator *validator.Validate
}

// Validate validates structs using go-playground/validator tags.
func (cv *customValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// NewHTTPServer creates and configures the Echo server with all middleware,
// dependency injection, and route registration.
func NewHTTPServer(cfg *config.Config, db *gorm.DB) *echo.Echo {
	e := echo.New()

	// ── Validator ─────────────────────────────────────────────────────────
	e.Validator = &customValidator{validator: validator.New()}

	// ── Global Middleware ──────────────────────────────────────────────────
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
		return c.JSON(http.StatusOK, httpresponse.Success{
			Success: true,
			Message: "SpotSync API is running 🚗",
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

	return e
}
