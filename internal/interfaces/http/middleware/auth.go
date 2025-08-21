package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// AuthMiddleware provides a simple authentication middleware
// TODO: Replace with proper JWT authentication
func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			if auth == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing authorization header"})
			}

			// Simple token validation (replace with JWT validation)
			if !strings.HasPrefix(auth, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid authorization header format"})
			}

			token := strings.TrimPrefix(auth, "Bearer ")
			if token == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing token"})
			}

			// TODO: Validate JWT token and extract user ID
			// For now, we'll use a placeholder
			c.Set("user_id", "placeholder-user-id")

			return next(c)
		}
	}
}
