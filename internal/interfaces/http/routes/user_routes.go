package routes

import (
	"nekosync/internal/domain/user"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterUserRoutes(e *echo.Echo, service user.Service) {
	userGroup := e.Group("/users")

	userGroup.POST("/", func(c echo.Context) error {
		var req struct {
			Email string `json:"email"`
			Name  string `json:"name"`
		}
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		u := &user.User{
			Email: req.Email,
			Name:  req.Name,
		}

		if err := service.RegisterUser(u); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, u)
	})

	userGroup.GET("/:id", func(c echo.Context) error {
		// dummy: add service.GetByID() when implemented
		return c.String(http.StatusOK, "Get user by ID not implemented")
	})
}
