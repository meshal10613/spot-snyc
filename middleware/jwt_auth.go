package middleware

import (
	"fmt"
	"net/http"
	"spot-sync/httpresponse"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// JWTAuth is a middleware that validates the Bearer token from the Authorization header,
// extracts user_id and role from the JWT payload, and injects them into Echo context.
func JWTAuth(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract Bearer token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, httpresponse.Error{
					Success: false,
					Message: "Missing authorization header",
				})
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, httpresponse.Error{
					Success: false,
					Message: "Invalid authorization header format. Use: Bearer <token>",
				})
			}

			tokenString := parts[1]

			// Parse and validate JWT
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, httpresponse.Error{
					Success: false,
					Message: "Invalid or expired token",
				})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, httpresponse.Error{
					Success: false,
					Message: "Invalid token claims",
				})
			}

			// Extract user_id (JWT encodes numbers as float64)
			userIDFloat, ok := claims["user_id"].(float64)
			if !ok {
				return c.JSON(http.StatusUnauthorized, httpresponse.Error{
					Success: false,
					Message: "Invalid user ID in token",
				})
			}

			role, ok := claims["role"].(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, httpresponse.Error{
					Success: false,
					Message: "Invalid role in token",
				})
			}

			// Inject user data into Echo context for downstream handlers
			c.Set("user_id", uint(userIDFloat))
			c.Set("role", role)

			return next(c)
		}
	}
}
