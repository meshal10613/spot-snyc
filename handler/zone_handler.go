package handler

import (
	"net/http"
	"spot-sync/dto"
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

	res, err := h.service.Create(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to create parking zone",
			"errors":  err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Parking zone created successfully",
		"data":    res,
	})
}

// GetAll handles GET /api/v1/zones (Public)
func (h *ZoneHandler) GetAll(c echo.Context) error {
	zones, err := h.service.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to retrieve parking zones",
			"errors":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Parking zones retrieved successfully",
		"data":    zones,
	})
}

// GetByID handles GET /api/v1/zones/:id (Public)
func (h *ZoneHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid zone ID",
		})
	}

	zone, err := h.service.GetByID(uint(id))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"success": false,
				"message": "Parking zone not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to retrieve parking zone",
			"errors":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Parking zone retrieved successfully",
		"data":    zone,
	})
}

// Update handles PUT /api/v1/zones/:id (Admin Only)
func (h *ZoneHandler) Update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid zone ID",
		})
	}

	var req dto.UpdateZoneRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid request body",
			"errors":  err.Error(),
		})
	}

	res, err := h.service.Update(uint(id), &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"success": false,
				"message": "Parking zone not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to update parking zone",
			"errors":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Parking zone updated successfully",
		"data":    res,
	})
}

// Delete handles DELETE /api/v1/zones/:id (Admin Only)
func (h *ZoneHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid zone ID",
		})
	}

	if err := h.service.Delete(uint(id)); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"success": false,
				"message": "Parking zone not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to delete parking zone",
			"errors":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Parking zone deleted successfully",
	})
}
