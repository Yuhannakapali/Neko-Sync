package app

import (
	"database/sql"
	"net/http"

	"nekosync/internal/platform/auth"
	"nekosync/internal/platform/config"
	"nekosync/internal/platform/httpx"
	"nekosync/internal/user"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// NewServer is the composition root: it builds every feature and registers its routes.
func NewServer(cfg *config.Config, db *sql.DB) *echo.Echo {
	server := echo.New()

	server.Use(middleware.Logger())
	server.Use(middleware.Recover())
	server.Use(middleware.CORS())

	server.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	tokens := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiry)

	api := server.Group("/api/v1")
	protected := api.Group("")
	protected.Use(httpx.AuthMiddleware(tokens))

	userService := user.NewService(
		user.NewUserRepository(db),
		user.NewDeviceRepository(db),
		user.NewFollowRepository(db),
		user.NewNotificationRepository(db),
	)
	user.NewHandler(userService, tokens).Routes(api, protected)

	return server
}
