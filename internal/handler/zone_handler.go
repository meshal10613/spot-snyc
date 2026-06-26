package handler

import (
	"net/http"
	"spot-sync/internal/dto"
	"spot-sync/pkg/httpresponse"
	"spot-sync/internal/service"
	"spot-sync/pkg/utils"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// ZoneHandler handles all parking zone HTTP endpoints.
type ZoneHandler struct {
	service service.ZoneService
}

// NewZoneHandler creates a new zone handler with injected service.
func NewZoneHandler(service service.ZoneService) *ZoneHandler {
	return &ZoneHandler{service: service}
}

// Create handles POST /api/v1/zones (Admin Only)
func (h *ZoneHandler) Create(c echo.Context) error {
	var req dto.CreateZoneRequest
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

	res, err := h.service.Create(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false,
			Message: "Failed to create parking zone",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, httpresponse.Success{
		Success: true,
		Message: "Parking zone created successfully",
		Data:    res,
	})
}

// GetAll handles GET /api/v1/zones (Public)
// Supports query params: page, limit, sort, order, search
func (h *ZoneHandler) GetAll(c echo.Context) error {
	qb := utils.NewQueryBuilder(c)

	zones, total, err := h.service.GetAll(qb)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false,
			Message: "Failed to retrieve parking zones",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true,
		Message: "Parking zones retrieved successfully",
		Data:    zones,
		Meta:    qb.GetMeta(total),
	})
}

// GetByID handles GET /api/v1/zones/:id (Public)
func (h *ZoneHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false,
			Message: "Invalid zone ID",
		})
	}

	zone, err := h.service.GetByID(uint(id))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, httpresponse.Error{
				Success: false,
				Message: "Parking zone not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false,
			Message: "Failed to retrieve parking zone",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true,
		Message: "Parking zone retrieved successfully",
		Data:    zone,
	})
}

// Update handles PUT /api/v1/zones/:id (Admin Only)
func (h *ZoneHandler) Update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false,
			Message: "Invalid zone ID",
		})
	}

	var req dto.UpdateZoneRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false,
			Message: "Invalid request body",
			Details: err.Error(),
		})
	}

	res, err := h.service.Update(uint(id), &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, httpresponse.Error{
				Success: false,
				Message: "Parking zone not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false,
			Message: "Failed to update parking zone",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true,
		Message: "Parking zone updated successfully",
		Data:    res,
	})
}

// Delete handles DELETE /api/v1/zones/:id (Admin Only)
func (h *ZoneHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false,
			Message: "Invalid zone ID",
		})
	}

	if err := h.service.Delete(uint(id)); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, httpresponse.Error{
				Success: false,
				Message: "Parking zone not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false,
			Message: "Failed to delete parking zone",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true,
		Message: "Parking zone deleted successfully",
	})
}
