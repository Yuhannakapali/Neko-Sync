package party

import (
	"context"
	"nekosync/internal/platform/entity"
)

type Repository interface {
	Create(ctx context.Context, p *WatchParty) error
	GetByID(ctx context.Context, id entity.UUID) (*WatchParty, error)
	GetByRoomCode(ctx context.Context, roomCode string) (*WatchParty, error)
	Update(ctx context.Context, p *WatchParty) error
	Delete(ctx context.Context, id entity.UUID) error
	GetActiveParties(ctx context.Context, limit, offset int) ([]*WatchParty, error)
	GetPartiesByHost(ctx context.Context, hostUserID entity.UUID) ([]*WatchParty, error)

	AddMember(ctx context.Context, m *PartyMember) error
	RemoveMember(ctx context.Context, partyID, userID entity.UUID) error
	GetMembers(ctx context.Context, partyID entity.UUID) ([]*PartyMember, error)
	GetUserParties(ctx context.Context, userID entity.UUID) ([]*WatchParty, error)
	IsMember(ctx context.Context, partyID, userID entity.UUID) (bool, error)
	GetMemberCount(ctx context.Context, partyID entity.UUID) (int, error)

	UpdatePlaybackState(ctx context.Context, state *PlaybackState) error
	GetPlaybackState(ctx context.Context, partyID entity.UUID) (*PlaybackState, error)

	CreateMessage(ctx context.Context, m *Message) error
	GetMessages(ctx context.Context, partyID entity.UUID, limit, offset int) ([]*Message, error)
	DeleteMessage(ctx context.Context, id entity.UUID) error
	GetMessagesByTimestamp(ctx context.Context, partyID entity.UUID, from, to int) ([]*Message, error)
}

type TransferRepository interface {
	Create(ctx context.Context, t *DeviceTransfer) error
	GetByID(ctx context.Context, id entity.UUID) (*DeviceTransfer, error)
	GetByUserID(ctx context.Context, userID entity.UUID) ([]*DeviceTransfer, error)
	Update(ctx context.Context, t *DeviceTransfer) error
	GetPending(ctx context.Context, userID entity.UUID) ([]*DeviceTransfer, error)
}
