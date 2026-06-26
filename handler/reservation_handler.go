package handler

import (
	"errors"
	"net/http"
	"spot-sync/dto"
	"spot-sync/httpresponse"
	"spot-sync/repository"
	"spot-sync/service"
	"spot-sync/utils"
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
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Success: false,
			Message: "Authentication required",
		})
	}

	var req dto.CreateReservationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false,
			Message: "Invalid request body",
			Details: err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false,
			Message: "Validation failed",
			Details: err.Error(),
		})
	}

	res, err := h.service.Create(userID, &req)
	if err != nil {
		// Zone full → 409 Conflict
		if errors.Is(err, repository.ErrZoneFull) {
			return c.JSON(http.StatusConflict, httpresponse.Error{
				Success: false,
				Message: "Parking zone is full. No available spots.",
			})
		}
		// Zone not found → 404
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, httpresponse.Error{
				Success: false,
				Message: "Parking zone not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false,
			Message: "Failed to create reservation",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, httpresponse.Success{
		Success: true,
		Message: "Reservation confirmed successfully",
		Data:    res,
	})
}

// GetMyReservations handles GET /api/v1/reservations/my-reservations (Authenticated)
// Supports query params: page, limit, sort, order, search
func (h *ReservationHandler) GetMyReservations(c echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Success: false,
			Message: "Authentication required",
		})
	}

	qb := utils.NewQueryBuilder(c)

	reservations, total, err := h.service.GetMyReservations(userID, qb)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false,
			Message: "Failed to retrieve reservations",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true,
		Message: "My reservations retrieved successfully",
		Data:    reservations,
		Meta:    qb.GetMeta(total),
	})
}

// Cancel handles DELETE /api/v1/reservations/:id (Authenticated)
func (h *ReservationHandler) Cancel(c echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Success: false,
			Message: "Authentication required",
		})
	}

	userRole, _ := c.Get("role").(string)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false,
			Message: "Invalid reservation ID",
		})
	}

	if err := h.service.Cancel(uint(id), userID, userRole); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, httpresponse.Error{
				Success: false,
				Message: "Reservation not found",
			})
		}
		if strings.Contains(err.Error(), "forbidden") {
			return c.JSON(http.StatusForbidden, httpresponse.Error{
				Success: false,
				Message: "You can only cancel your own reservations",
			})
		}
		if strings.Contains(err.Error(), "already cancelled") {
			return c.JSON(http.StatusBadRequest, httpresponse.Error{
				Success: false,
				Message: "Reservation is already cancelled",
			})
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false,
			Message: "Failed to cancel reservation",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true,
		Message: "Reservation cancelled successfully",
	})
}

// GetAll handles GET /api/v1/reservations (Admin Only)
// Supports query params: page, limit, sort, order, search
func (h *ReservationHandler) GetAll(c echo.Context) error {
	qb := utils.NewQueryBuilder(c)

	reservations, total, err := h.service.GetAll(qb)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false,
			Message: "Failed to retrieve reservations",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true,
		Message: "All reservations retrieved successfully",
		Data:    reservations,
		Meta:    qb.GetMeta(total),
	})
}
