package repositories

import (
	"context"
	"database/sql"
	"nekosync/internal/domain/shared"
	"nekosync/internal/domain/user"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) user.Repository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *user.User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash, avatar_url, role, is_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.ExecContext(ctx, query,
		u.ID, u.Username, u.Email, u.PasswordHash,
		u.AvatarURL, u.Role, u.IsVerified,
		u.CreatedAt, u.UpdatedAt)
	return err
}

func (r *userRepository) GetByID(ctx context.Context, id shared.UUID) (*user.User, error) {
	u := &user.User{}
	query := `
		SELECT id, username, email, password_hash, avatar_url, role, is_verified, created_at, updated_at
		FROM users WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash,
		&u.AvatarURL, &u.Role, &u.IsVerified,
		&u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	u := &user.User{}
	query := `
		SELECT id, username, email, password_hash, avatar_url, role, is_verified, created_at, updated_at
		FROM users WHERE email = $1`

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash,
		&u.AvatarURL, &u.Role, &u.IsVerified,
		&u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	u := &user.User{}
	query := `
		SELECT id, username, email, password_hash, avatar_url, role, is_verified, created_at, updated_at
		FROM users WHERE username = $1`

	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash,
		&u.AvatarURL, &u.Role, &u.IsVerified,
		&u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepository) Update(ctx context.Context, u *user.User) error {
	query := `
		UPDATE users SET username = $2, email = $3, password_hash = $4, avatar_url = $5,
		role = $6, is_verified = $7, updated_at = $8 WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		u.ID, u.Username, u.Email, u.PasswordHash,
		u.AvatarURL, u.Role, u.IsVerified, u.UpdatedAt)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id shared.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	return err
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*user.User, error) {
	query := `
		SELECT id, username, email, password_hash, avatar_url, role, is_verified, created_at, updated_at
		FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*user.User
	for rows.Next() {
		u := &user.User{}
		if err := rows.Scan(
			&u.ID, &u.Username, &u.Email, &u.PasswordHash,
			&u.AvatarURL, &u.Role, &u.IsVerified,
			&u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *userRepository) CreateProfile(ctx context.Context, p *user.Profile) error {
	query := `
		INSERT INTO user_profiles (user_id, about, location, website, banner_url, birthdate)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.ExecContext(ctx, query,
		p.UserID, p.About, p.Location, p.Website, p.BannerURL, p.Birthdate)
	return err
}

func (r *userRepository) GetProfile(ctx context.Context, userID shared.UUID) (*user.Profile, error) {
	p := &user.Profile{}
	query := `
		SELECT user_id, about, location, website, banner_url, birthdate
		FROM user_profiles WHERE user_id = $1`

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&p.UserID, &p.About, &p.Location, &p.Website, &p.BannerURL, &p.Birthdate)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *userRepository) UpdateProfile(ctx context.Context, p *user.Profile) error {
	query := `
		UPDATE user_profiles SET about = $2, location = $3, website = $4,
		banner_url = $5, birthdate = $6 WHERE user_id = $1`

	_, err := r.db.ExecContext(ctx, query,
		p.UserID, p.About, p.Location, p.Website, p.BannerURL, p.Birthdate)
	return err
}
