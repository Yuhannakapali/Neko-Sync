package handlers

import (
	"nekosync/internal/application/dto"
	"nekosync/internal/application/usecases/user"
	"net/http"

	"github.com/labstack/echo/v4"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	createUserUseCase       *user.CreateUserUseCase
	authenticateUserUseCase *user.AuthenticateUserUseCase
	updateProfileUseCase    *user.UpdateProfileUseCase
	followUserUseCase       *user.FollowUserUseCase
	registerDeviceUseCase   *user.RegisterDeviceUseCase
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(
	createUserUseCase *user.CreateUserUseCase,
	authenticateUserUseCase *user.AuthenticateUserUseCase,
	updateProfileUseCase *user.UpdateProfileUseCase,
	followUserUseCase *user.FollowUserUseCase,
	registerDeviceUseCase *user.RegisterDeviceUseCase,
) *UserHandler {
	return &UserHandler{
		createUserUseCase:       createUserUseCase,
		authenticateUserUseCase: authenticateUserUseCase,
		updateProfileUseCase:    updateProfileUseCase,
		followUserUseCase:       followUserUseCase,
		registerDeviceUseCase:   registerDeviceUseCase,
	}
}

// CreateUser handles user registration
func (h *UserHandler) CreateUser(c echo.Context) error {
	var req dto.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// TODO: Add validation

	resp, err := h.createUserUseCase.Execute(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, resp)
}

// Login handles user authentication
func (h *UserHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// TODO: Add validation

	resp, err := h.authenticateUserUseCase.Execute(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

// UpdateProfile handles user profile updates
func (h *UserHandler) UpdateProfile(c echo.Context) error {
	// TODO: Extract user ID from JWT token
	userID := c.Get("user_id").(string)

	var req dto.UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	resp, err := h.updateProfileUseCase.Execute(c.Request().Context(), userID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

// FollowUser handles user following
func (h *UserHandler) FollowUser(c echo.Context) error {
	// TODO: Extract user ID from JWT token
	followerID := c.Get("user_id").(string)

	var req dto.FollowUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	err := h.followUserUseCase.Execute(c.Request().Context(), followerID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Successfully followed user"})
}

// RegisterDevice handles device registration
func (h *UserHandler) RegisterDevice(c echo.Context) error {
	// TODO: Extract user ID from JWT token
	userID := c.Get("user_id").(string)

	var req dto.RegisterDeviceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	resp, err := h.registerDeviceUseCase.Execute(c.Request().Context(), userID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, resp)
}
