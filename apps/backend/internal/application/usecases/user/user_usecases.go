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
	userService services.UserServiceInterface
}

func NewCreateUserUseCase(userService services.UserServiceInterface) *CreateUserUseCase {
	return &CreateUserUseCase{userService: userService}
}

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
	userService services.UserServiceInterface
}

func NewAuthenticateUserUseCase(userService services.UserServiceInterface) *AuthenticateUserUseCase {
	return &AuthenticateUserUseCase{userService: userService}
}

func (uc *AuthenticateUserUseCase) Execute(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := uc.userService.AuthenticateUser(ctx, req.Email, req.Password)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// TODO: Generate real JWT token
	token := "jwt-token-placeholder"

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
	userService services.UserServiceInterface
}

func NewUpdateProfileUseCase(userService services.UserServiceInterface) *UpdateProfileUseCase {
	return &UpdateProfileUseCase{userService: userService}
}

func (uc *UpdateProfileUseCase) Execute(ctx context.Context, userID string, req dto.UpdateProfileRequest) (*dto.UserProfileResponse, error) {
	profile := &entities.UserProfile{
		UserID:    entities.UUID(userID),
		About:     req.About,
		Location:  req.Location,
		Website:   req.Website,
		BannerURL: req.BannerURL,
		Birthdate: req.Birthdate,
	}

	if err := uc.userService.UpdateProfile(ctx, entities.UUID(userID), profile); err != nil {
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
	userService services.UserServiceInterface
}

func NewFollowUserUseCase(userService services.UserServiceInterface) *FollowUserUseCase {
	return &FollowUserUseCase{userService: userService}
}

func (uc *FollowUserUseCase) Execute(ctx context.Context, followerID string, req dto.FollowUserRequest) error {
	if err := uc.userService.FollowUser(ctx, entities.UUID(followerID), entities.UUID(req.UserID)); err != nil {
		return fmt.Errorf("failed to follow user: %w", err)
	}

	// Fire-and-forget notification; errors are intentionally not surfaced to the caller.
	uc.userService.CreateNotification(
		ctx,
		entities.UUID(req.UserID),
		entities.NotificationTypeFollow,
		"New Follower",
		"Someone started following you",
		map[string]interface{}{"follower_id": followerID},
	)

	return nil
}

// RegisterDeviceUseCase handles device registration.
type RegisterDeviceUseCase struct {
	userService services.UserServiceInterface
}

func NewRegisterDeviceUseCase(userService services.UserServiceInterface) *RegisterDeviceUseCase {
	return &RegisterDeviceUseCase{userService: userService}
}

func (uc *RegisterDeviceUseCase) Execute(ctx context.Context, userID string, req dto.RegisterDeviceRequest) (*dto.DeviceResponse, error) {
	device, err := uc.userService.RegisterDevice(ctx, entities.UUID(userID), req.DeviceName, entities.PlatformType(req.Platform))
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
