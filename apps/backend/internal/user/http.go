package user

import (
	"nekosync/internal/platform/entity"
	"net/http"

	"github.com/labstack/echo/v4"
)

const timeFormat = "2006-01-02T15:04:05Z"

type TokenIssuer interface {
	Issue(userID string) (string, error)
}

// Handler handles HTTP requests for user operations
type Handler struct {
	svc    *Service
	tokens TokenIssuer
}

func NewHandler(svc *Service, tokens TokenIssuer) *Handler {
	return &Handler{svc: svc, tokens: tokens}
}

// CreateUser handles user registration
func (h *Handler) CreateUser(c echo.Context) error {
	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// TODO: Add validation

	u, err := h.svc.CreateUser(c.Request().Context(), req.Username, req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "failed to create user: " + err.Error()})
	}

	return c.JSON(http.StatusCreated, CreateUserResponse{
		ID:       string(u.ID),
		Username: u.Username,
		Email:    u.Email,
		Role:     string(u.Role),
	})
}

// Login handles user authentication
func (h *Handler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// TODO: Add validation

	u, err := h.svc.AuthenticateUser(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "authentication failed: " + err.Error()})
	}

	token, err := h.tokens.Issue(string(u.ID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to issue token"})
	}

	return c.JSON(http.StatusOK, LoginResponse{
		User: UserResponse{
			ID:         string(u.ID),
			Username:   u.Username,
			Email:      u.Email,
			AvatarURL:  u.AvatarURL,
			Role:       string(u.Role),
			IsVerified: u.IsVerified,
			CreatedAt:  u.CreatedAt.Format(timeFormat),
		},
		Token: token,
	})
}

// UpdateProfile handles user profile updates
func (h *Handler) UpdateProfile(c echo.Context) error {
	userID := c.Get("user_id").(string)

	var req UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	profile := &Profile{
		UserID:    entity.UUID(userID),
		About:     req.About,
		Location:  req.Location,
		Website:   req.Website,
		BannerURL: req.BannerURL,
		Birthdate: req.Birthdate,
	}

	if err := h.svc.UpdateProfile(c.Request().Context(), entity.UUID(userID), profile); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "failed to update profile: " + err.Error()})
	}

	return c.JSON(http.StatusOK, UserProfileResponse{
		UserID:    userID,
		About:     profile.About,
		Location:  profile.Location,
		Website:   profile.Website,
		BannerURL: profile.BannerURL,
		Birthdate: profile.Birthdate,
	})
}

// FollowUser handles user following
func (h *Handler) FollowUser(c echo.Context) error {
	followerID := c.Get("user_id").(string)

	var req FollowUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	ctx := c.Request().Context()
	if err := h.svc.FollowUser(ctx, entity.UUID(followerID), entity.UUID(req.UserID)); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "failed to follow user: " + err.Error()})
	}

	// Fire-and-forget notification; errors are intentionally not surfaced to the caller.
	h.svc.CreateNotification(
		ctx,
		entity.UUID(req.UserID),
		NotificationFollow,
		"New Follower",
		"Someone started following you",
		map[string]interface{}{"follower_id": followerID},
	)

	return c.JSON(http.StatusOK, map[string]string{"message": "Successfully followed user"})
}

// RegisterDevice handles device registration
func (h *Handler) RegisterDevice(c echo.Context) error {
	userID := c.Get("user_id").(string)

	var req RegisterDeviceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	device, err := h.svc.RegisterDevice(c.Request().Context(), entity.UUID(userID), req.DeviceName, PlatformType(req.Platform))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "failed to register device: " + err.Error()})
	}

	return c.JSON(http.StatusCreated, DeviceResponse{
		ID:         string(device.ID),
		DeviceName: device.DeviceName,
		Platform:   string(device.Platform),
		LastSeen:   device.LastSeen.Format(timeFormat),
		IsActive:   device.IsActive,
	})
}

// Routes registers the user endpoints. protected must already require authentication.
func (h *Handler) Routes(public, protected *echo.Group) {
	public.POST("/users/register", h.CreateUser)
	public.POST("/users/login", h.Login)

	protected.PUT("/users/profile", h.UpdateProfile)
	protected.POST("/users/follow", h.FollowUser)
	protected.POST("/users/devices", h.RegisterDevice)
}
