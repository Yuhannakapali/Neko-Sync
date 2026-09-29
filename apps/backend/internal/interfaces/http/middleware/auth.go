package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type TokenVerifier interface {
	Verify(token string) (userID string, err error)
}

// AuthMiddleware requires a valid bearer token and stores its subject as "user_id".
func AuthMiddleware(verifier TokenVerifier) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			if auth == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing authorization header"})
			}

			if !strings.HasPrefix(auth, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid authorization header format"})
			}

			token := strings.TrimPrefix(auth, "Bearer ")
			if token == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing token"})
			}

			userID, err := verifier.Verify(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid or expired token"})
			}

			c.Set("user_id", userID)

			return next(c)
		}
	}
}
