package user

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"nekosync/internal/domain/shared"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// ServiceInterface is the contract used by the application layer.
type ServiceInterface interface {
	CreateUser(ctx context.Context, username, email, password string) (*User, error)
	AuthenticateUser(ctx context.Context, email, password string) (*User, error)
	UpdateProfile(ctx context.Context, userID shared.UUID, profile *Profile) error
	FollowUser(ctx context.Context, followerID, followingID shared.UUID) error
	UnfollowUser(ctx context.Context, followerID, followingID shared.UUID) error
	CreateNotification(ctx context.Context, userID shared.UUID, notifType NotificationType, title, message string, data map[string]interface{}) error
	RegisterDevice(ctx context.Context, userID shared.UUID, deviceName string, platform PlatformType) (*Device, error)
}

type Service struct {
	repo         Repository
	deviceRepo   DeviceRepository
	followRepo   FollowRepository
	notifRepo    NotificationRepository
}

func NewService(
	repo Repository,
	deviceRepo DeviceRepository,
	followRepo FollowRepository,
	notifRepo NotificationRepository,
) *Service {
	return &Service{
		repo:       repo,
		deviceRepo: deviceRepo,
		followRepo: followRepo,
		notifRepo:  notifRepo,
	}
}

func (s *Service) CreateUser(ctx context.Context, username, email, password string) (*User, error) {
	if username == "" || email == "" || password == "" {
		return nil, fmt.Errorf("username, email, and password are required")
	}

	if existing, _ := s.repo.GetByEmail(ctx, email); existing != nil {
		return nil, ErrAlreadyExists
	}

	if existing, _ := s.repo.GetByUsername(ctx, username); existing != nil {
		return nil, ErrAlreadyExists
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	u := &User{
		BaseEntity: shared.BaseEntity{
			ID:        shared.UUID(generateID()),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Username:     username,
		Email:        email,
		PasswordHash: string(hashed),
		Role:         RoleUser,
		IsVerified:   false,
	}

	if err := s.repo.Create(ctx, u); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return u, nil
}

func (s *Service) AuthenticateUser(ctx context.Context, email, password string) (*User, error) {
	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return nil, ErrInvalidCredentials
	}

	return u, nil
}

func (s *Service) UpdateProfile(ctx context.Context, userID shared.UUID, profile *Profile) error {
	u, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return ErrNotFound
	}

	profile.UserID = u.ID

	existing, _ := s.repo.GetProfile(ctx, userID)
	if existing == nil {
		return s.repo.CreateProfile(ctx, profile)
	}

	return s.repo.UpdateProfile(ctx, profile)
}

func (s *Service) FollowUser(ctx context.Context, followerID, followingID shared.UUID) error {
	if followerID == followingID {
		return ErrCannotFollowSelf
	}

	if _, err := s.repo.GetByID(ctx, followerID); err != nil {
		return ErrNotFound
	}

	if _, err := s.repo.GetByID(ctx, followingID); err != nil {
		return ErrNotFound
	}

	isFollowing, err := s.followRepo.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return fmt.Errorf("failed to check follow status: %w", err)
	}

	if isFollowing {
		return ErrAlreadyFollowing
	}

	return s.followRepo.Follow(ctx, followerID, followingID)
}

func (s *Service) UnfollowUser(ctx context.Context, followerID, followingID shared.UUID) error {
	isFollowing, err := s.followRepo.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return fmt.Errorf("failed to check follow status: %w", err)
	}

	if !isFollowing {
		return ErrNotFollowing
	}

	return s.followRepo.Unfollow(ctx, followerID, followingID)
}

func (s *Service) CreateNotification(ctx context.Context, userID shared.UUID, notifType NotificationType, title, message string, data map[string]interface{}) error {
	n := &Notification{
		BaseEntity: shared.BaseEntity{
			ID:        shared.UUID(generateID()),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserID:  userID,
		Type:    notifType,
		Title:   title,
		Message: message,
		IsRead:  false,
		Data:    data,
	}

	return s.notifRepo.Create(ctx, n)
}

func (s *Service) RegisterDevice(ctx context.Context, userID shared.UUID, deviceName string, platform PlatformType) (*Device, error) {
	if err := s.deviceRepo.DeactivateAllForUser(ctx, userID); err != nil {
		return nil, fmt.Errorf("failed to deactivate existing devices: %w", err)
	}

	d := &Device{
		BaseEntity: shared.BaseEntity{
			ID:        shared.UUID(generateID()),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserID:     userID,
		DeviceName: deviceName,
		Platform:   platform,
		LastSeen:   time.Now(),
		IsActive:   true,
	}

	if err := s.deviceRepo.Create(ctx, d); err != nil {
		return nil, fmt.Errorf("failed to create device: %w", err)
	}

	return d, nil
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
