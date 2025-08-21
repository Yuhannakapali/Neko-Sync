package party

import "nekosync/internal/domain"

type Repository interface {
	// Watch Party operations
	CreateWatchParty(party *WatchParty) error
	GetWatchParty(id domain.UUID) (*WatchParty, error)
	GetWatchPartyByRoomCode(roomCode string) (*WatchParty, error)
	GetActiveWatchParties() ([]*WatchParty, error)
	GetUserWatchParties(userID domain.UUID) ([]*WatchParty, error)
	UpdateWatchParty(party *WatchParty) error
	EndWatchParty(id domain.UUID) error
	DeleteWatchParty(id domain.UUID) error

	// Watch Party User operations
	AddWatchPartyUser(partyUser *WatchPartyUser) error
	RemoveWatchPartyUser(partyID, userID domain.UUID) error
	GetWatchPartyUsers(partyID domain.UUID) ([]*WatchPartyUser, error)
	IsUserInParty(partyID, userID domain.UUID) (bool, error)

	// Device Transfer operations
	CreateDeviceTransfer(transfer *DeviceTransfer) error
	GetDeviceTransfer(id domain.UUID) (*DeviceTransfer, error)
	GetDeviceTransfersByUser(userID domain.UUID) ([]*DeviceTransfer, error)
	GetDeviceTransfersByDevice(deviceID domain.UUID) ([]*DeviceTransfer, error)
	DeleteDeviceTransfer(id domain.UUID) error
}
