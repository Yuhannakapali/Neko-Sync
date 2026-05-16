package user

import (
	"context"
	"fmt"
	"nekosync/internal/application/dto"
	"nekosync/internal/domain/shared"
	userDomain "nekosync/internal/domain/user"
)

type CreateUserUseCase struct {
	svc userDomain.ServiceInterface
}

func NewCreateUserUseCase(svc userDomain.ServiceInterface) *CreateUserUseCase {
	return &CreateUserUseCase{svc: svc}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, req dto.CreateUserRequest) (*dto.CreateUserResponse, error) {
	u, err := uc.svc.CreateUser(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &dto.CreateUserResponse{
		ID:       string(u.ID),
		Username: u.Username,
		Email:    u.Email,
		Role:     string(u.Role),
	}, nil
}

type AuthenticateUserUseCase struct {
	svc userDomain.ServiceInterface
}

func NewAuthenticateUserUseCase(svc userDomain.ServiceInterface) *AuthenticateUserUseCase {
	return &AuthenticateUserUseCase{svc: svc}
}

func (uc *AuthenticateUserUseCase) Execute(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	u, err := uc.svc.AuthenticateUser(ctx, req.Email, req.Password)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// TODO: Generate real JWT token
	token := "jwt-token-placeholder"

	return &dto.LoginResponse{
		User: dto.UserResponse{
			ID:         string(u.ID),
			Username:   u.Username,
			Email:      u.Email,
			AvatarURL:  u.AvatarURL,
			Role:       string(u.Role),
			IsVerified: u.IsVerified,
			CreatedAt:  u.CreatedAt.Format("2006-01-02T15:04:05Z"),
		},
		Token: token,
	}, nil
}

type UpdateProfileUseCase struct {
	svc userDomain.ServiceInterface
}

func NewUpdateProfileUseCase(svc userDomain.ServiceInterface) *UpdateProfileUseCase {
	return &UpdateProfileUseCase{svc: svc}
}

func (uc *UpdateProfileUseCase) Execute(ctx context.Context, userID string, req dto.UpdateProfileRequest) (*dto.UserProfileResponse, error) {
	profile := &userDomain.Profile{
		UserID:    shared.UUID(userID),
		About:     req.About,
		Location:  req.Location,
		Website:   req.Website,
		BannerURL: req.BannerURL,
		Birthdate: req.Birthdate,
	}

	if err := uc.svc.UpdateProfile(ctx, shared.UUID(userID), profile); err != nil {
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

type FollowUserUseCase struct {
	svc userDomain.ServiceInterface
}

func NewFollowUserUseCase(svc userDomain.ServiceInterface) *FollowUserUseCase {
	return &FollowUserUseCase{svc: svc}
}

func (uc *FollowUserUseCase) Execute(ctx context.Context, followerID string, req dto.FollowUserRequest) error {
	if err := uc.svc.FollowUser(ctx, shared.UUID(followerID), shared.UUID(req.UserID)); err != nil {
		return fmt.Errorf("failed to follow user: %w", err)
	}

	// Fire-and-forget notification; errors are intentionally not surfaced to the caller.
	uc.svc.CreateNotification(
		ctx,
		shared.UUID(req.UserID),
		userDomain.NotificationFollow,
		"New Follower",
		"Someone started following you",
		map[string]interface{}{"follower_id": followerID},
	)

	return nil
}

type RegisterDeviceUseCase struct {
	svc userDomain.ServiceInterface
}

func NewRegisterDeviceUseCase(svc userDomain.ServiceInterface) *RegisterDeviceUseCase {
	return &RegisterDeviceUseCase{svc: svc}
}

func (uc *RegisterDeviceUseCase) Execute(ctx context.Context, userID string, req dto.RegisterDeviceRequest) (*dto.DeviceResponse, error) {
	device, err := uc.svc.RegisterDevice(ctx, shared.UUID(userID), req.DeviceName, userDomain.PlatformType(req.Platform))
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
