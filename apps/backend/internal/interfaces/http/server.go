package http

import (
	"database/sql"
	"nekosync/internal/application/usecases/user"
	"nekosync/internal/config"
	userDomain "nekosync/internal/domain/user"
	"nekosync/internal/infrastructure/repositories"
	"nekosync/internal/interfaces/http/handlers"
	customMiddleware "nekosync/internal/interfaces/http/middleware"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewHTTPServer(cfg *config.Config, db *sql.DB) *echo.Echo {
	server := echo.New()

	server.Use(middleware.Logger())
	server.Use(middleware.Recover())
	server.Use(middleware.CORS())

	server.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// Repositories
	userRepo := repositories.NewUserRepository(db)
	deviceRepo := repositories.NewDeviceRepository(db)
	followRepo := repositories.NewFollowRepository(db)
	notifRepo := repositories.NewNotificationRepository(db)

	// Domain service
	userService := userDomain.NewService(userRepo, deviceRepo, followRepo, notifRepo)

	// Use cases
	createUserUC := user.NewCreateUserUseCase(userService)
	authenticateUC := user.NewAuthenticateUserUseCase(userService)
	updateProfileUC := user.NewUpdateProfileUseCase(userService)
	followUserUC := user.NewFollowUserUseCase(userService)
	registerDeviceUC := user.NewRegisterDeviceUseCase(userService)

	// Handlers
	userHandler := handlers.NewUserHandler(
		createUserUC,
		authenticateUC,
		updateProfileUC,
		followUserUC,
		registerDeviceUC,
	)

	// Routes
	api := server.Group("/api/v1")

	api.POST("/users/register", userHandler.CreateUser)
	api.POST("/users/login", userHandler.Login)

	protected := api.Group("")
	protected.Use(customMiddleware.AuthMiddleware())

	protected.PUT("/users/profile", userHandler.UpdateProfile)
	protected.POST("/users/follow", userHandler.FollowUser)
	protected.POST("/users/devices", userHandler.RegisterDevice)

	return server
}
