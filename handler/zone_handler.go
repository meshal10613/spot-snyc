package handler

import (
	"net/http"
	"spot-sync/dto"
	"spot-sync/httpresponse"
	"spot-sync/service"
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
func (h *ZoneHandler) GetAll(c echo.Context) error {
	zones, err := h.service.GetAll()
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
