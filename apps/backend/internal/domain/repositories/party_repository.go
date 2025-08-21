package repositories

import (
	"context"
	"nekosync/internal/domain/entities"
)

// PartyRepository defines the contract for watch party data operations
type PartyRepository interface {
	// Watch Party CRUD operations
	Create(ctx context.Context, party *entities.WatchParty) error
	GetByID(ctx context.Context, id entities.UUID) (*entities.WatchParty, error)
	GetByRoomCode(ctx context.Context, roomCode string) (*entities.WatchParty, error)
	Update(ctx context.Context, party *entities.WatchParty) error
	Delete(ctx context.Context, id entities.UUID) error
	GetActiveParties(ctx context.Context, limit, offset int) ([]*entities.WatchParty, error)
	GetPartiesByHost(ctx context.Context, hostUserID entities.UUID) ([]*entities.WatchParty, error)

	// Watch Party User operations
	AddUser(ctx context.Context, partyUser *entities.WatchPartyUser) error
	RemoveUser(ctx context.Context, partyID, userID entities.UUID) error
	GetPartyUsers(ctx context.Context, partyID entities.UUID) ([]*entities.WatchPartyUser, error)
	GetUserParties(ctx context.Context, userID entities.UUID) ([]*entities.WatchParty, error)
	IsUserInParty(ctx context.Context, partyID, userID entities.UUID) (bool, error)
	GetPartyUserCount(ctx context.Context, partyID entities.UUID) (int, error)

	// Watch Party State operations
	UpdateState(ctx context.Context, state *entities.WatchPartyState) error
	GetState(ctx context.Context, partyID entities.UUID) (*entities.WatchPartyState, error)

	// Device Transfer operations
	CreateTransfer(ctx context.Context, transfer *entities.DeviceTransfer) error
	GetTransferByID(ctx context.Context, id entities.UUID) (*entities.DeviceTransfer, error)
	GetTransfersByUser(ctx context.Context, userID entities.UUID) ([]*entities.DeviceTransfer, error)
	UpdateTransfer(ctx context.Context, transfer *entities.DeviceTransfer) error
	GetPendingTransfers(ctx context.Context, userID entities.UUID) ([]*entities.DeviceTransfer, error)

	// Party Chat operations
	CreateMessage(ctx context.Context, message *entities.PartyMessage) error
	GetMessages(ctx context.Context, partyID entities.UUID, limit, offset int) ([]*entities.PartyMessage, error)
	DeleteMessage(ctx context.Context, id entities.UUID) error
	GetMessagesByTimestamp(ctx context.Context, partyID entities.UUID, fromTimestamp, toTimestamp int) ([]*entities.PartyMessage, error)
}
