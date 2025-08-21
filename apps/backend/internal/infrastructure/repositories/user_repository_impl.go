package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"nekosync/internal/domain/entities"
	"nekosync/internal/domain/repositories"
	"time"
)

// UserRepositoryImpl implements the UserRepository interface
type UserRepositoryImpl struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepositoryImpl
func NewUserRepository(db *sql.DB) repositories.UserRepository {
	return &UserRepositoryImpl{db: db}
}

// User CRUD operations
func (r *UserRepositoryImpl) Create(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash, avatar_url, role, is_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Username, user.Email, user.PasswordHash,
		user.AvatarURL, user.Role, user.IsVerified,
		user.CreatedAt, user.UpdatedAt)

	return err
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, id entities.UUID) (*entities.User, error) {
	user := &entities.User{}
	query := `
		SELECT id, username, email, password_hash, avatar_url, role, is_verified, created_at, updated_at
		FROM users WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.AvatarURL, &user.Role, &user.IsVerified,
		&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	user := &entities.User{}
	query := `
		SELECT id, username, email, password_hash, avatar_url, role, is_verified, created_at, updated_at
		FROM users WHERE email = $1`

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.AvatarURL, &user.Role, &user.IsVerified,
		&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepositoryImpl) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
	user := &entities.User{}
	query := `
		SELECT id, username, email, password_hash, avatar_url, role, is_verified, created_at, updated_at
		FROM users WHERE username = $1`

	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.AvatarURL, &user.Role, &user.IsVerified,
		&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, user *entities.User) error {
	query := `
		UPDATE users SET username = $2, email = $3, password_hash = $4, avatar_url = $5, 
		role = $6, is_verified = $7, updated_at = $8
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Username, user.Email, user.PasswordHash,
		user.AvatarURL, user.Role, user.IsVerified, user.UpdatedAt)

	return err
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id entities.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *UserRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	query := `
		SELECT id, username, email, password_hash, avatar_url, role, is_verified, created_at, updated_at
		FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		user := &entities.User{}
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.PasswordHash,
			&user.AvatarURL, &user.Role, &user.IsVerified,
			&user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// User Profile operations
func (r *UserRepositoryImpl) CreateProfile(ctx context.Context, profile *entities.UserProfile) error {
	query := `
		INSERT INTO user_profiles (user_id, about, location, website, banner_url, birthdate)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.ExecContext(ctx, query,
		profile.UserID, profile.About, profile.Location,
		profile.Website, profile.BannerURL, profile.Birthdate)

	return err
}

func (r *UserRepositoryImpl) GetProfile(ctx context.Context, userID entities.UUID) (*entities.UserProfile, error) {
	profile := &entities.UserProfile{}
	query := `
		SELECT user_id, about, location, website, banner_url, birthdate
		FROM user_profiles WHERE user_id = $1`

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&profile.UserID, &profile.About, &profile.Location,
		&profile.Website, &profile.BannerURL, &profile.Birthdate)

	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (r *UserRepositoryImpl) UpdateProfile(ctx context.Context, profile *entities.UserProfile) error {
	query := `
		UPDATE user_profiles SET about = $2, location = $3, website = $4, 
		banner_url = $5, birthdate = $6
		WHERE user_id = $1`

	_, err := r.db.ExecContext(ctx, query,
		profile.UserID, profile.About, profile.Location,
		profile.Website, profile.BannerURL, profile.Birthdate)

	return err
}

// User Device operations
func (r *UserRepositoryImpl) CreateDevice(ctx context.Context, device *entities.UserDevice) error {
	query := `
		INSERT INTO user_devices (id, user_id, device_name, platform, last_seen, websocket_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.ExecContext(ctx, query,
		device.ID, device.UserID, device.DeviceName, device.Platform,
		device.LastSeen, device.WebsocketID, device.IsActive,
		device.CreatedAt, device.UpdatedAt)

	return err
}

func (r *UserRepositoryImpl) GetDevicesByUserID(ctx context.Context, userID entities.UUID) ([]*entities.UserDevice, error) {
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
		err := rows.Scan(
			&device.ID, &device.UserID, &device.DeviceName, &device.Platform,
			&device.LastSeen, &device.WebsocketID, &device.IsActive,
			&device.CreatedAt, &device.UpdatedAt)
		if err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}

	return devices, nil
}

func (r *UserRepositoryImpl) UpdateDevice(ctx context.Context, device *entities.UserDevice) error {
	query := `
		UPDATE user_devices SET device_name = $2, platform = $3, last_seen = $4, 
		websocket_id = $5, is_active = $6, updated_at = $7
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		device.ID, device.DeviceName, device.Platform,
		device.LastSeen, device.WebsocketID, device.IsActive, device.UpdatedAt)

	return err
}

func (r *UserRepositoryImpl) DeleteDevice(ctx context.Context, id entities.UUID) error {
	query := `DELETE FROM user_devices WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *UserRepositoryImpl) DeactivateUserDevices(ctx context.Context, userID entities.UUID) error {
	query := `UPDATE user_devices SET is_active = false, updated_at = $2 WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID, time.Now())
	return err
}

// User Follow operations
func (r *UserRepositoryImpl) Follow(ctx context.Context, followerID, followingID entities.UUID) error {
	query := `INSERT INTO user_follows (follower_id, following_id, followed_at) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, followerID, followingID, time.Now())
	return err
}

func (r *UserRepositoryImpl) Unfollow(ctx context.Context, followerID, followingID entities.UUID) error {
	query := `DELETE FROM user_follows WHERE follower_id = $1 AND following_id = $2`
	_, err := r.db.ExecContext(ctx, query, followerID, followingID)
	return err
}

func (r *UserRepositoryImpl) GetFollowers(ctx context.Context, userID entities.UUID, limit, offset int) ([]*entities.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.avatar_url, u.role, u.is_verified, u.created_at, u.updated_at
		FROM users u
		INNER JOIN user_follows uf ON u.id = uf.follower_id
		WHERE uf.following_id = $1
		ORDER BY uf.followed_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		user := &entities.User{}
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.PasswordHash,
			&user.AvatarURL, &user.Role, &user.IsVerified,
			&user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepositoryImpl) GetFollowing(ctx context.Context, userID entities.UUID, limit, offset int) ([]*entities.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.avatar_url, u.role, u.is_verified, u.created_at, u.updated_at
		FROM users u
		INNER JOIN user_follows uf ON u.id = uf.following_id
		WHERE uf.follower_id = $1
		ORDER BY uf.followed_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		user := &entities.User{}
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.PasswordHash,
			&user.AvatarURL, &user.Role, &user.IsVerified,
			&user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepositoryImpl) IsFollowing(ctx context.Context, followerID, followingID entities.UUID) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM user_follows WHERE follower_id = $1 AND following_id = $2`
	err := r.db.QueryRowContext(ctx, query, followerID, followingID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepositoryImpl) GetFollowCount(ctx context.Context, userID entities.UUID) (followers int, following int, err error) {
	// Get followers count
	query := `SELECT COUNT(*) FROM user_follows WHERE following_id = $1`
	err = r.db.QueryRowContext(ctx, query, userID).Scan(&followers)
	if err != nil {
		return 0, 0, err
	}

	// Get following count
	query = `SELECT COUNT(*) FROM user_follows WHERE follower_id = $1`
	err = r.db.QueryRowContext(ctx, query, userID).Scan(&following)
	if err != nil {
		return 0, 0, err
	}

	return followers, following, nil
}

// Notification operations
func (r *UserRepositoryImpl) CreateNotification(ctx context.Context, notification *entities.Notification) error {
	dataJSON, _ := json.Marshal(notification.Data)

	query := `
		INSERT INTO notifications (id, user_id, type, title, message, is_read, data, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.ExecContext(ctx, query,
		notification.ID, notification.UserID, notification.Type,
		notification.Title, notification.Message, notification.IsRead,
		dataJSON, notification.CreatedAt, notification.UpdatedAt)

	return err
}

func (r *UserRepositoryImpl) GetNotifications(ctx context.Context, userID entities.UUID, limit, offset int) ([]*entities.Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, is_read, data, created_at, updated_at
		FROM notifications WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*entities.Notification
	for rows.Next() {
		notification := &entities.Notification{}
		var dataJSON []byte

		err := rows.Scan(
			&notification.ID, &notification.UserID, &notification.Type,
			&notification.Title, &notification.Message, &notification.IsRead,
			&dataJSON, &notification.CreatedAt, &notification.UpdatedAt)
		if err != nil {
			return nil, err
		}

		if len(dataJSON) > 0 {
			json.Unmarshal(dataJSON, &notification.Data)
		}

		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (r *UserRepositoryImpl) GetUnreadNotifications(ctx context.Context, userID entities.UUID) ([]*entities.Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, is_read, data, created_at, updated_at
		FROM notifications WHERE user_id = $1 AND is_read = false
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*entities.Notification
	for rows.Next() {
		notification := &entities.Notification{}
		var dataJSON []byte

		err := rows.Scan(
			&notification.ID, &notification.UserID, &notification.Type,
			&notification.Title, &notification.Message, &notification.IsRead,
			&dataJSON, &notification.CreatedAt, &notification.UpdatedAt)
		if err != nil {
			return nil, err
		}

		if len(dataJSON) > 0 {
			json.Unmarshal(dataJSON, &notification.Data)
		}

		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (r *UserRepositoryImpl) MarkNotificationAsRead(ctx context.Context, notificationID entities.UUID) error {
	query := `UPDATE notifications SET is_read = true, updated_at = $2 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, notificationID, time.Now())
	return err
}

func (r *UserRepositoryImpl) MarkAllNotificationsAsRead(ctx context.Context, userID entities.UUID) error {
	query := `UPDATE notifications SET is_read = true, updated_at = $2 WHERE user_id = $1 AND is_read = false`
	_, err := r.db.ExecContext(ctx, query, userID, time.Now())
	return err
}

func (r *UserRepositoryImpl) DeleteNotification(ctx context.Context, id entities.UUID) error {
	query := `DELETE FROM notifications WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
