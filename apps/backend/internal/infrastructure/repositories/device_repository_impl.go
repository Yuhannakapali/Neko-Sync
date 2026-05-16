package repositories

import (
	"context"
	"database/sql"
	"nekosync/internal/domain/shared"
	"nekosync/internal/domain/user"
	"time"
)

type deviceRepository struct {
	db *sql.DB
}

func NewDeviceRepository(db *sql.DB) user.DeviceRepository {
	return &deviceRepository{db: db}
}

func (r *deviceRepository) Create(ctx context.Context, d *user.Device) error {
	query := `
		INSERT INTO user_devices (id, user_id, device_name, platform, last_seen, websocket_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.ExecContext(ctx, query,
		d.ID, d.UserID, d.DeviceName, d.Platform,
		d.LastSeen, d.WebsocketID, d.IsActive,
		d.CreatedAt, d.UpdatedAt)
	return err
}

func (r *deviceRepository) GetByUserID(ctx context.Context, userID shared.UUID) ([]*user.Device, error) {
	query := `
		SELECT id, user_id, device_name, platform, last_seen, websocket_id, is_active, created_at, updated_at
		FROM user_devices WHERE user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []*user.Device
	for rows.Next() {
		d := &user.Device{}
		if err := rows.Scan(
			&d.ID, &d.UserID, &d.DeviceName, &d.Platform,
			&d.LastSeen, &d.WebsocketID, &d.IsActive,
			&d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		devices = append(devices, d)
	}
	return devices, nil
}

func (r *deviceRepository) Update(ctx context.Context, d *user.Device) error {
	query := `
		UPDATE user_devices SET device_name = $2, platform = $3, last_seen = $4,
		websocket_id = $5, is_active = $6, updated_at = $7 WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		d.ID, d.DeviceName, d.Platform,
		d.LastSeen, d.WebsocketID, d.IsActive, d.UpdatedAt)
	return err
}

func (r *deviceRepository) Delete(ctx context.Context, id shared.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM user_devices WHERE id = $1`, id)
	return err
}

func (r *deviceRepository) DeactivateAllForUser(ctx context.Context, userID shared.UUID) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE user_devices SET is_active = false, updated_at = $2 WHERE user_id = $1`,
		userID, time.Now())
	return err
}
