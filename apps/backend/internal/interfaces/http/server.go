package http

import (
	"database/sql"
	"nekosync/internal/application/usecases/user"
	"nekosync/internal/config"
	"nekosync/internal/domain/services"
	"nekosync/internal/infrastructure/repositories"
	"nekosync/internal/interfaces/http/handlers"
	customMiddleware "nekosync/internal/interfaces/http/middleware"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewHTTPServer(cfg *config.Config, db *sql.DB) *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Health check route
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)

	// Initialize domain services
	userService := services.NewUserService(userRepo)

	// Initialize use cases
	createUserUseCase := user.NewCreateUserUseCase(userService)
	authenticateUserUseCase := user.NewAuthenticateUserUseCase(userService)
	updateProfileUseCase := user.NewUpdateProfileUseCase(userService)
	followUserUseCase := user.NewFollowUserUseCase(userService)
	registerDeviceUseCase := user.NewRegisterDeviceUseCase(userService)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(
		createUserUseCase,
		authenticateUserUseCase,
		updateProfileUseCase,
		followUserUseCase,
		registerDeviceUseCase,
	)

	// Routes
	api := e.Group("/api/v1")

	// User routes
	api.POST("/users/register", userHandler.CreateUser)
	api.POST("/users/login", userHandler.Login)

	// Protected routes (require authentication middleware)
	protected := api.Group("")
	protected.Use(customMiddleware.AuthMiddleware())

	protected.PUT("/users/profile", userHandler.UpdateProfile)
	protected.POST("/users/follow", userHandler.FollowUser)
	protected.POST("/users/devices", userHandler.RegisterDevice)

	return e
}
