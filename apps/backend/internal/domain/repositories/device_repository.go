package repositories

import (
	"context"
	"nekosync/internal/domain/entities"
)

type DeviceRepository interface {
	Create(ctx context.Context, device *entities.UserDevice) error
	GetByUserID(ctx context.Context, userID entities.UUID) ([]*entities.UserDevice, error)
	Update(ctx context.Context, device *entities.UserDevice) error
	Delete(ctx context.Context, id entities.UUID) error
	DeactivateAllForUser(ctx context.Context, userID entities.UUID) error
}
