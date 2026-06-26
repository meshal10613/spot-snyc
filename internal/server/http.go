package server

import (
	"net/http"
	"spot-sync/internal/config"
	"spot-sync/internal/handler"
	"spot-sync/pkg/httpresponse"
	"spot-sync/internal/repository"
	"spot-sync/internal/routes"
	"spot-sync/internal/service"

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

// globalErrorHandler is a centralized error handler that catches all unhandled
// errors (404, 405, panics, validation errors, etc.) and returns a consistent
// httpresponse.Error JSON response.
func globalErrorHandler(err error, c echo.Context) {
	// Don't overwrite a response that was already committed
	if c.Response().Committed {
		return
	}

	code := http.StatusInternalServerError
	message := "Internal server error"
	details := ""

	// Echo wraps route-not-found, method-not-allowed, etc. as *echo.HTTPError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code

		switch code {
		case http.StatusNotFound:
			message = "The requested resource was not found"
		case http.StatusMethodNotAllowed:
			message = "Method not allowed"
		default:
			// Use the message from the HTTPError if available
			if msg, ok := he.Message.(string); ok {
				message = msg
			}
		}
	} else {
		// Non-HTTP errors (unexpected panics, DB failures, etc.)
		details = err.Error()
	}

	_ = c.JSON(code, httpresponse.Error{
		Success: false,
		Message: message,
		Details: details,
	})
}

// NewHTTPServer creates and configures the Echo server with all middleware,
// dependency injection, and route registration.
func NewHTTPServer(cfg *config.Config, db *gorm.DB) *echo.Echo {
	e := echo.New()

	// ── Global Error Handler ──────────────────────────────────────────────
	e.HTTPErrorHandler = globalErrorHandler

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
