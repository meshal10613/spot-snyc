package handler

import (
	"errors"
	"net/http"
	"spot-sync/dto"
	"spot-sync/repository"
	"spot-sync/service"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// ReservationHandler handles all reservation HTTP endpoints.
type ReservationHandler struct {
	service service.ReservationService
}

// NewReservationHandler creates a new reservation handler with injected service.
func NewReservationHandler(service service.ReservationService) *ReservationHandler {
	return &ReservationHandler{service: service}
}

// Create handles POST /api/v1/reservations (Authenticated)
func (h *ReservationHandler) Create(c echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "Authentication required",
		})
	}

	var req dto.CreateReservationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid request body",
			"errors":  err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Validation failed",
			"errors":  err.Error(),
		})
	}

	res, err := h.service.Create(userID, &req)
	if err != nil {
		// Zone full → 409 Conflict
		if errors.Is(err, repository.ErrZoneFull) {
			return c.JSON(http.StatusConflict, map[string]interface{}{
				"success": false,
				"message": "Parking zone is full. No available spots.",
			})
		}
		// Zone not found → 404
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"success": false,
				"message": "Parking zone not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to create reservation",
			"errors":  err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Reservation confirmed successfully",
		"data":    res,
	})
}

// GetMyReservations handles GET /api/v1/reservations/my-reservations (Authenticated)
func (h *ReservationHandler) GetMyReservations(c echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "Authentication required",
		})
	}

	reservations, err := h.service.GetMyReservations(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to retrieve reservations",
			"errors":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "My reservations retrieved successfully",
		"data":    reservations,
	})
}

// Cancel handles DELETE /api/v1/reservations/:id (Authenticated)
func (h *ReservationHandler) Cancel(c echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "Authentication required",
		})
	}

	userRole, _ := c.Get("role").(string)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid reservation ID",
		})
	}

	if err := h.service.Cancel(uint(id), userID, userRole); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"success": false,
				"message": "Reservation not found",
			})
		}
		if strings.Contains(err.Error(), "forbidden") {
			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"success": false,
				"message": "You can only cancel your own reservations",
			})
		}
		if strings.Contains(err.Error(), "already cancelled") {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "Reservation is already cancelled",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to cancel reservation",
			"errors":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Reservation cancelled successfully",
	})
}

// GetAll handles GET /api/v1/reservations (Admin Only)
func (h *ReservationHandler) GetAll(c echo.Context) error {
	reservations, err := h.service.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to retrieve reservations",
			"errors":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "All reservations retrieved successfully",
		"data":    reservations,
	})
}
