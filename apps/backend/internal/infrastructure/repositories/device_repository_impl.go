package repositories

import (
	"context"
	"database/sql"
	"nekosync/internal/domain/entities"
	"nekosync/internal/domain/repositories"
	"time"
)

type deviceRepository struct {
	db *sql.DB
}

func NewDeviceRepository(db *sql.DB) repositories.DeviceRepository {
	return &deviceRepository{db: db}
}

func (r *deviceRepository) Create(ctx context.Context, device *entities.UserDevice) error {
	query := `
		INSERT INTO user_devices (id, user_id, device_name, platform, last_seen, websocket_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.ExecContext(ctx, query,
		device.ID, device.UserID, device.DeviceName, device.Platform,
		device.LastSeen, device.WebsocketID, device.IsActive,
		device.CreatedAt, device.UpdatedAt)
	return err
}

func (r *deviceRepository) GetByUserID(ctx context.Context, userID entities.UUID) ([]*entities.UserDevice, error) {
	query := `
		SELECT id, user_id, device_name, platform, last_seen, websocket_id, is_active, created_at, updated_at
		FROM user_devices WHERE user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []*entities.UserDevice
	for rows.Next() {
		device := &entities.UserDevice{}
		if err := rows.Scan(
			&device.ID, &device.UserID, &device.DeviceName, &device.Platform,
			&device.LastSeen, &device.WebsocketID, &device.IsActive,
			&device.CreatedAt, &device.UpdatedAt); err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}
	return devices, nil
}

func (r *deviceRepository) Update(ctx context.Context, device *entities.UserDevice) error {
	query := `
		UPDATE user_devices SET device_name = $2, platform = $3, last_seen = $4,
		websocket_id = $5, is_active = $6, updated_at = $7
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		device.ID, device.DeviceName, device.Platform,
		device.LastSeen, device.WebsocketID, device.IsActive, device.UpdatedAt)
	return err
}

func (r *deviceRepository) Delete(ctx context.Context, id entities.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM user_devices WHERE id = $1`, id)
	return err
}

func (r *deviceRepository) DeactivateAllForUser(ctx context.Context, userID entities.UUID) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE user_devices SET is_active = false, updated_at = $2 WHERE user_id = $1`,
		userID, time.Now())
	return err
}
