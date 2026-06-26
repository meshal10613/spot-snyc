package handler

import (
	"net/http"
	"spot-sync/internal/dto"
	"spot-sync/pkg/httpresponse"
	"spot-sync/internal/service"
	"strings"

	"github.com/labstack/echo/v4"
)

// AuthHandler handles all authentication HTTP endpoints.
type AuthHandler struct {
	service service.AuthService
}

// NewAuthHandler creates a new auth handler with injected service.
func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Register handles POST /api/v1/auth/register
func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
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

	res, err := h.service.Register(&req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return c.JSON(http.StatusBadRequest, httpresponse.Error{
				Success: false,
				Message: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false,
			Message: "Registration failed",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, httpresponse.Success{
		Success: true,
		Message: "User registered successfully",
		Data:    res,
	})
}

// Login handles POST /api/v1/auth/login
func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
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

	res, err := h.service.Login(&req)
	if err != nil {
		if strings.Contains(err.Error(), "invalid email or password") {
			return c.JSON(http.StatusUnauthorized, httpresponse.Error{
				Success: false,
				Message: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false,
			Message: "Login failed",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true,
		Message: "Login successful",
		Data:    res,
	})
}
