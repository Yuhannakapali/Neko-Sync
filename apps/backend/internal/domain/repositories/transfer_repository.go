package repositories

import (
	"context"
	"nekosync/internal/domain/entities"
)

type DeviceTransferRepository interface {
	Create(ctx context.Context, transfer *entities.DeviceTransfer) error
	GetByID(ctx context.Context, id entities.UUID) (*entities.DeviceTransfer, error)
	GetByUserID(ctx context.Context, userID entities.UUID) ([]*entities.DeviceTransfer, error)
	Update(ctx context.Context, transfer *entities.DeviceTransfer) error
	GetPending(ctx context.Context, userID entities.UUID) ([]*entities.DeviceTransfer, error)
}
