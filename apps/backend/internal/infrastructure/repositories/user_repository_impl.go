package repositories

import (
	"context"
	"database/sql"
	"nekosync/internal/domain/entities"
	"nekosync/internal/domain/repositories"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repositories.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash, avatar_url, role, is_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Username, user.Email, user.PasswordHash,
		user.AvatarURL, user.Role, user.IsVerified,
		user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *userRepository) GetByID(ctx context.Context, id entities.UUID) (*entities.User, error) {
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

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
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

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
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

func (r *userRepository) Update(ctx context.Context, user *entities.User) error {
	query := `
		UPDATE users SET username = $2, email = $3, password_hash = $4, avatar_url = $5,
		role = $6, is_verified = $7, updated_at = $8
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Username, user.Email, user.PasswordHash,
		user.AvatarURL, user.Role, user.IsVerified, user.UpdatedAt)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id entities.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	return err
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*entities.User, error) {
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
		if err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.PasswordHash,
			&user.AvatarURL, &user.Role, &user.IsVerified,
			&user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepository) CreateProfile(ctx context.Context, profile *entities.UserProfile) error {
	query := `
		INSERT INTO user_profiles (user_id, about, location, website, banner_url, birthdate)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.ExecContext(ctx, query,
		profile.UserID, profile.About, profile.Location,
		profile.Website, profile.BannerURL, profile.Birthdate)
	return err
}

func (r *userRepository) GetProfile(ctx context.Context, userID entities.UUID) (*entities.UserProfile, error) {
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

func (r *userRepository) UpdateProfile(ctx context.Context, profile *entities.UserProfile) error {
	query := `
		UPDATE user_profiles SET about = $2, location = $3, website = $4,
		banner_url = $5, birthdate = $6
		WHERE user_id = $1`

	_, err := r.db.ExecContext(ctx, query,
		profile.UserID, profile.About, profile.Location,
		profile.Website, profile.BannerURL, profile.Birthdate)
	return err
}
