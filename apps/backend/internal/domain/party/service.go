package party

import (
	"context"
	"crypto/rand"
	"fmt"
	"nekosync/internal/domain/content"
	"nekosync/internal/domain/shared"
	"nekosync/internal/domain/user"
	"time"
)

// ServiceInterface is the contract used by the application layer.
type ServiceInterface interface {
	CreateWatchParty(ctx context.Context, hostUserID, contentID shared.UUID, title string, maxUsers int, isPrivate bool, password *string) (*WatchParty, error)
	JoinWatchParty(ctx context.Context, userID shared.UUID, roomCode string, password *string) (*WatchParty, error)
	LeaveWatchParty(ctx context.Context, userID, partyID shared.UUID) error
	UpdatePlaybackState(ctx context.Context, partyID, userID shared.UUID, currentTime int, isPlaying bool, speed float64) error
	SendMessage(ctx context.Context, partyID, userID shared.UUID, message string, timestamp int) error
	CreateDeviceTransfer(ctx context.Context, userID, fromDeviceID, toDeviceID, contentID shared.UUID, position int) (*DeviceTransfer, error)
}

type Service struct {
	partyRepo    Repository
	transferRepo TransferRepository
	contentRepo  content.Repository
	userRepo     user.Repository
	deviceRepo   user.DeviceRepository
}

func NewService(
	partyRepo Repository,
	transferRepo TransferRepository,
	contentRepo content.Repository,
	userRepo user.Repository,
	deviceRepo user.DeviceRepository,
) *Service {
	return &Service{
		partyRepo:    partyRepo,
		transferRepo: transferRepo,
		contentRepo:  contentRepo,
		userRepo:     userRepo,
		deviceRepo:   deviceRepo,
	}
}

func (s *Service) CreateWatchParty(ctx context.Context, hostUserID, contentID shared.UUID, title string, maxUsers int, isPrivate bool, password *string) (*WatchParty, error) {
	if _, err := s.userRepo.GetByID(ctx, hostUserID); err != nil {
		return nil, user.ErrNotFound
	}

	if _, err := s.contentRepo.GetByID(ctx, contentID); err != nil {
		return nil, content.ErrNotFound
	}

	p := &WatchParty{
		BaseEntity: shared.BaseEntity{
			ID:        shared.UUID(generateID()),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		HostUserID: hostUserID,
		ContentID:  contentID,
		RoomCode:   generateRoomCode(),
		Title:      title,
		MaxUsers:   maxUsers,
		IsPrivate:  isPrivate,
		Password:   password,
		IsActive:   true,
	}

	if err := s.partyRepo.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("failed to create watch party: %w", err)
	}

	if err := s.partyRepo.AddMember(ctx, &PartyMember{
		PartyID:  p.ID,
		UserID:   hostUserID,
		JoinedAt: time.Now(),
		IsActive: true,
	}); err != nil {
		return nil, fmt.Errorf("failed to add host to party: %w", err)
	}

	return p, nil
}

func (s *Service) JoinWatchParty(ctx context.Context, userID shared.UUID, roomCode string, password *string) (*WatchParty, error) {
	p, err := s.partyRepo.GetByRoomCode(ctx, roomCode)
	if err != nil {
		return nil, ErrNotFound
	}

	if !p.IsActive {
		return nil, ErrNotActive
	}

	if p.IsPrivate {
		if password == nil || p.Password == nil || *password != *p.Password {
			return nil, ErrWrongPassword
		}
	}

	isMember, err := s.partyRepo.IsMember(ctx, p.ID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check party membership: %w", err)
	}

	if isMember {
		return nil, ErrAlreadyMember
	}

	count, err := s.partyRepo.GetMemberCount(ctx, p.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get party member count: %w", err)
	}

	if p.IsFull(count) {
		return nil, ErrFull
	}

	if err := s.partyRepo.AddMember(ctx, &PartyMember{
		PartyID:  p.ID,
		UserID:   userID,
		JoinedAt: time.Now(),
		IsActive: true,
	}); err != nil {
		return nil, fmt.Errorf("failed to join party: %w", err)
	}

	return p, nil
}

func (s *Service) LeaveWatchParty(ctx context.Context, userID, partyID shared.UUID) error {
	p, err := s.partyRepo.GetByID(ctx, partyID)
	if err != nil {
		return ErrNotFound
	}

	isMember, err := s.partyRepo.IsMember(ctx, partyID, userID)
	if err != nil {
		return fmt.Errorf("failed to check party membership: %w", err)
	}

	if !isMember {
		return ErrNotMember
	}

	if err := s.partyRepo.RemoveMember(ctx, partyID, userID); err != nil {
		return fmt.Errorf("failed to leave party: %w", err)
	}

	if p.IsHost(userID) {
		p.End()
		if err := s.partyRepo.Update(ctx, p); err != nil {
			return fmt.Errorf("failed to end party: %w", err)
		}
	}

	return nil
}

func (s *Service) UpdatePlaybackState(ctx context.Context, partyID, userID shared.UUID, currentTime int, isPlaying bool, speed float64) error {
	p, err := s.partyRepo.GetByID(ctx, partyID)
	if err != nil {
		return ErrNotFound
	}

	if !p.IsActive {
		return ErrNotActive
	}

	isMember, err := s.partyRepo.IsMember(ctx, partyID, userID)
	if err != nil {
		return fmt.Errorf("failed to check party membership: %w", err)
	}

	if !isMember {
		return ErrNotMember
	}

	return s.partyRepo.UpdatePlaybackState(ctx, &PlaybackState{
		PartyID:       partyID,
		CurrentTime:   currentTime,
		IsPlaying:     isPlaying,
		PlaybackSpeed: speed,
		UpdatedAt:     time.Now(),
		UpdatedBy:     userID,
	})
}

func (s *Service) SendMessage(ctx context.Context, partyID, userID shared.UUID, message string, timestamp int) error {
	isMember, err := s.partyRepo.IsMember(ctx, partyID, userID)
	if err != nil {
		return fmt.Errorf("failed to check party membership: %w", err)
	}

	if !isMember {
		return ErrNotMember
	}

	return s.partyRepo.CreateMessage(ctx, &Message{
		BaseEntity: shared.BaseEntity{
			ID:        shared.UUID(generateID()),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		PartyID:   partyID,
		UserID:    userID,
		Message:   message,
		Timestamp: timestamp,
	})
}

func (s *Service) CreateDeviceTransfer(ctx context.Context, userID, fromDeviceID, toDeviceID, contentID shared.UUID, position int) (*DeviceTransfer, error) {
	devices, err := s.deviceRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user devices: %w", err)
	}

	fromExists, toExists := false, false
	for _, d := range devices {
		if d.ID == fromDeviceID {
			fromExists = true
		}
		if d.ID == toDeviceID {
			toExists = true
		}
	}

	if !fromExists || !toExists {
		return nil, ErrInvalidDevices
	}

	c, err := s.contentRepo.GetByID(ctx, contentID)
	if err != nil {
		return nil, content.ErrNotFound
	}

	t := &DeviceTransfer{
		BaseEntity: shared.BaseEntity{
			ID:        shared.UUID(generateID()),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserID:       userID,
		FromDeviceID: fromDeviceID,
		ToDeviceID:   toDeviceID,
		ContentType:  c.Type,
		ContentID:    contentID,
		Position:     position,
		IsCompleted:  false,
	}

	if err := s.transferRepo.Create(ctx, t); err != nil {
		return nil, fmt.Errorf("failed to create device transfer: %w", err)
	}

	return t, nil
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	result := make([]byte, 32)
	const hex = "0123456789abcdef"
	for i, v := range b {
		result[i*2] = hex[v>>4]
		result[i*2+1] = hex[v&0x0f]
	}
	return string(result)
}

func generateRoomCode() string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)
	for i := range b {
		r := make([]byte, 1)
		rand.Read(r)
		b[i] = chars[int(r[0])%len(chars)]
	}
	return string(b)
}
