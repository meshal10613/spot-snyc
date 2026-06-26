package routes

import (
	"spot-sync/internal/handler"
	"spot-sync/pkg/middleware"
	"spot-sync/internal/models"

	"github.com/labstack/echo/v4"
)

// RegisterRoutes wires all API endpoints with their handlers and middleware.
func RegisterRoutes(
	e *echo.Echo,
	authHandler *handler.AuthHandler,
	zoneHandler *handler.ZoneHandler,
	reservationHandler *handler.ReservationHandler,
	jwtSecret string,
) {
	api := e.Group("/api/v1")

	// ── Authentication (Public) ────────────────────────────────────────────
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	// ── Parking Zones ──────────────────────────────────────────────────────
	zones := api.Group("/zones")
	// Public: anyone can view zones
	zones.GET("", zoneHandler.GetAll)
	zones.GET("/:id", zoneHandler.GetByID)
	// Admin only: create, update, delete (per-route middleware)
	zones.POST("", zoneHandler.Create, middleware.JWTAuth(jwtSecret), middleware.RequireRole(string(models.RoleAdmin)))
	zones.PUT("/:id", zoneHandler.Update, middleware.JWTAuth(jwtSecret), middleware.RequireRole(string(models.RoleAdmin)))
	zones.DELETE("/:id", zoneHandler.Delete, middleware.JWTAuth(jwtSecret), middleware.RequireRole(string(models.RoleAdmin)))

	// ── Reservations (Authenticated) ───────────────────────────────────────
	reservations := api.Group("/reservations", middleware.JWTAuth(jwtSecret))
	reservations.POST("", reservationHandler.Create)
	reservations.GET("/my-reservations", reservationHandler.GetMyReservations)
	reservations.DELETE("/:id", reservationHandler.Cancel)
	// Admin only: view all reservations (additional role check)
	reservations.GET("", reservationHandler.GetAll, middleware.RequireRole(string(models.RoleAdmin)))
}

