package http

import (
	"database/sql"
	"nekosync/internal/config"
	"nekosync/internal/domain/user"
	"nekosync/internal/infra/persistence"
	"nekosync/internal/interfaces/http/routes"
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewServer(cfg *config.Config, db *sql.DB) *echo.Echo {
	e := echo.New()

	// Basic route
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// Initialize domain service/repo
	userRepo := persistence.NewUserRepo(db)
	userService := user.NewService(userRepo)

	// Register routes per domain
	routes.RegisterUserRoutes(e, *userService)

	return e
}
