package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// RequireRole is a middleware that checks if the authenticated user has one of the allowed roles.
// Must be used after JWTAuth middleware which injects "role" into context.
func RequireRole(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, ok := c.Get("role").(string)
			if !ok {
				return c.JSON(http.StatusForbidden, map[string]interface{}{
					"success": false,
					"message": "Access denied",
				})
			}

			for _, allowedRole := range roles {
				if role == allowedRole {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"success": false,
				"message": "You don't have permission to access this resource",
			})
		}
	}
}
