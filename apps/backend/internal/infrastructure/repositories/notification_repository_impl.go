package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"nekosync/internal/domain/entities"
	"nekosync/internal/domain/repositories"
	"time"
)

type notificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) repositories.NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, n *entities.Notification) error {
	dataJSON, _ := json.Marshal(n.Data)

	query := `
		INSERT INTO notifications (id, user_id, type, title, message, is_read, data, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.ExecContext(ctx, query,
		n.ID, n.UserID, n.Type, n.Title, n.Message, n.IsRead,
		dataJSON, n.CreatedAt, n.UpdatedAt)
	return err
}

func (r *notificationRepository) GetByUserID(ctx context.Context, userID entities.UUID, limit, offset int) ([]*entities.Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, is_read, data, created_at, updated_at
		FROM notifications WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	return r.scanNotifications(ctx, query, userID, limit, offset)
}

func (r *notificationRepository) GetUnread(ctx context.Context, userID entities.UUID) ([]*entities.Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, is_read, data, created_at, updated_at
		FROM notifications WHERE user_id = $1 AND is_read = false
		ORDER BY created_at DESC`

	return r.scanNotifications(ctx, query, userID)
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, notificationID entities.UUID) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE notifications SET is_read = true, updated_at = $2 WHERE id = $1`,
		notificationID, time.Now())
	return err
}

func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID entities.UUID) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE notifications SET is_read = true, updated_at = $2 WHERE user_id = $1 AND is_read = false`,
		userID, time.Now())
	return err
}

func (r *notificationRepository) Delete(ctx context.Context, id entities.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM notifications WHERE id = $1`, id)
	return err
}

func (r *notificationRepository) scanNotifications(ctx context.Context, query string, args ...interface{}) ([]*entities.Notification, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*entities.Notification
	for rows.Next() {
		n := &entities.Notification{}
		var dataJSON []byte

		if err := rows.Scan(
			&n.ID, &n.UserID, &n.Type, &n.Title, &n.Message, &n.IsRead,
			&dataJSON, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}

		if len(dataJSON) > 0 {
			json.Unmarshal(dataJSON, &n.Data)
		}

		notifications = append(notifications, n)
	}
	return notifications, nil
}
