package user

import (
	"context"
	"fmt"

	"nekosync/internal/application/dto"
	"nekosync/internal/domain/entities"
	"nekosync/internal/domain/services"
)

// CreateUserUseCase handles user creation.
type CreateUserUseCase struct {
	userService *services.UserService
}

// NewCreateUserUseCase creates a new CreateUserUseCase.
func NewCreateUserUseCase(userService *services.UserService) *CreateUserUseCase {
	return &CreateUserUseCase{
		userService: userService,
	}
}

// Execute creates a new user.
func (uc *CreateUserUseCase) Execute(ctx context.Context, req dto.CreateUserRequest) (*dto.CreateUserResponse, error) {
	user, err := uc.userService.CreateUser(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &dto.CreateUserResponse{
		ID:       string(user.ID),
		Username: user.Username,
		Email:    user.Email,
		Role:     string(user.Role),
	}, nil
}

// AuthenticateUserUseCase handles user authentication.
type AuthenticateUserUseCase struct {
	userService *services.UserService
}

// NewAuthenticateUserUseCase creates a new AuthenticateUserUseCase.
func NewAuthenticateUserUseCase(userService *services.UserService) *AuthenticateUserUseCase {
	return &AuthenticateUserUseCase{
		userService: userService,
	}
}

// Execute authenticates a user.
func (uc *AuthenticateUserUseCase) Execute(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := uc.userService.AuthenticateUser(ctx, req.Email, req.Password)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Here you would generate a JWT token
	token := "jwt-token-placeholder" // TODO: Implement JWT token generation

	return &dto.LoginResponse{
		User: dto.UserResponse{
			ID:         string(user.ID),
			Username:   user.Username,
			Email:      user.Email,
			AvatarURL:  user.AvatarURL,
			Role:       string(user.Role),
			IsVerified: user.IsVerified,
			CreatedAt:  user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		},
		Token: token,
	}, nil
}

// UpdateProfileUseCase handles user profile updates.
type UpdateProfileUseCase struct {
	userService *services.UserService
}

// NewUpdateProfileUseCase creates a new UpdateProfileUseCase.
func NewUpdateProfileUseCase(userService *services.UserService) *UpdateProfileUseCase {
	return &UpdateProfileUseCase{
		userService: userService,
	}
}

// Execute updates a user's profile.
func (uc *UpdateProfileUseCase) Execute(ctx context.Context, userID string, req dto.UpdateProfileRequest) (*dto.UserProfileResponse, error) {
	profile := &entities.UserProfile{
		UserID:    entities.UUID(userID),
		About:     req.About,
		Location:  req.Location,
		Website:   req.Website,
		BannerURL: req.BannerURL,
		Birthdate: req.Birthdate,
	}

	err := uc.userService.UpdateProfile(ctx, entities.UUID(userID), profile)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	return &dto.UserProfileResponse{
		UserID:    userID,
		About:     profile.About,
		Location:  profile.Location,
		Website:   profile.Website,
		BannerURL: profile.BannerURL,
		Birthdate: profile.Birthdate,
	}, nil
}

// FollowUserUseCase handles user following.
type FollowUserUseCase struct {
	userService *services.UserService
}

// NewFollowUserUseCase creates a new FollowUserUseCase.
func NewFollowUserUseCase(userService *services.UserService) *FollowUserUseCase {
	return &FollowUserUseCase{
		userService: userService,
	}
}

// Execute follows a user.
func (uc *FollowUserUseCase) Execute(ctx context.Context, followerID string, req dto.FollowUserRequest) error {
	err := uc.userService.FollowUser(ctx, entities.UUID(followerID), entities.UUID(req.UserID))
	if err != nil {
		return fmt.Errorf("failed to follow user: %w", err)
	}

	// Create follow notification
	err = uc.userService.CreateNotification(
		ctx,
		entities.UUID(req.UserID),
		entities.NotificationTypeFollow,
		"New Follower",
		"Someone started following you",
		map[string]interface{}{
			"follower_id": followerID,
		},
	)
	if err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return nil
}

// RegisterDeviceUseCase handles device registration.
type RegisterDeviceUseCase struct {
	userService *services.UserService
}

// NewRegisterDeviceUseCase creates a new RegisterDeviceUseCase.
func NewRegisterDeviceUseCase(userService *services.UserService) *RegisterDeviceUseCase {
	return &RegisterDeviceUseCase{
		userService: userService,
	}
}

// Execute registers a new device for a user.
func (uc *RegisterDeviceUseCase) Execute(ctx context.Context, userID string, req dto.RegisterDeviceRequest) (*dto.DeviceResponse, error) {
	platform := entities.PlatformType(req.Platform)

	device, err := uc.userService.RegisterDevice(ctx, entities.UUID(userID), req.DeviceName, platform)
	if err != nil {
		return nil, fmt.Errorf("failed to register device: %w", err)
	}

	return &dto.DeviceResponse{
		ID:         string(device.ID),
		DeviceName: device.DeviceName,
		Platform:   string(device.Platform),
		LastSeen:   device.LastSeen.Format("2006-01-02T15:04:05Z"),
		IsActive:   device.IsActive,
	}, nil
}
