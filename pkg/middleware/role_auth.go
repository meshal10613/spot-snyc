package middleware

import (
	"net/http"
	"spot-sync/pkg/httpresponse"

	"github.com/labstack/echo/v4"
)

// RequireRole is a middleware that checks if the authenticated user has one of the allowed roles.
// Must be used after JWTAuth middleware which injects "role" into context.
func RequireRole(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, ok := c.Get("role").(string)
			if !ok {
				return c.JSON(http.StatusForbidden, httpresponse.Error{
					Success: false,
					Message: "Access denied",
				})
			}

			for _, allowedRole := range roles {
				if role == allowedRole {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, httpresponse.Error{
				Success: false,
				Message: "You don't have permission to access this resource",
			})
		}
	}
}
