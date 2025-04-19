package http

import (
	"database/sql"
	"nekosync/internal/config"
	"nekosync/internal/domain/user"
	"nekosync/internal/infra/persistence"
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewServer(cfg *config.Config, db *sql.DB) *echo.Echo {
	e := echo.New()

	repo := persistence.NewUserRepo(db)
	service := user.NewService(repo)

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	e.POST("/users", func(c echo.Context) error {
		type req struct {
			Email string `json:"email"`
			Name  string `json:"name"`
		}
		var r req
		if err := c.Bind(&r); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		u := &user.User{Email: r.Email, Name: r.Name}
		err := service.RegisterUser(u)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, u)
	})

	return e
}
