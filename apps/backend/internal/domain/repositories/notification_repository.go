package repositories

import (
	"context"
	"nekosync/internal/domain/entities"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *entities.Notification) error
	GetByUserID(ctx context.Context, userID entities.UUID, limit, offset int) ([]*entities.Notification, error)
	GetUnread(ctx context.Context, userID entities.UUID) ([]*entities.Notification, error)
	MarkAsRead(ctx context.Context, notificationID entities.UUID) error
	MarkAllAsRead(ctx context.Context, userID entities.UUID) error
	Delete(ctx context.Context, id entities.UUID) error
}
