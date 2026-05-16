package repositories

import (
	"context"
	"database/sql"
	"nekosync/internal/domain/entities"
	"nekosync/internal/domain/repositories"
	"time"
)

type followRepository struct {
	db *sql.DB
}

func NewFollowRepository(db *sql.DB) repositories.FollowRepository {
	return &followRepository{db: db}
}

func (r *followRepository) Follow(ctx context.Context, followerID, followingID entities.UUID) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO user_follows (follower_id, following_id, followed_at) VALUES ($1, $2, $3)`,
		followerID, followingID, time.Now())
	return err
}

func (r *followRepository) Unfollow(ctx context.Context, followerID, followingID entities.UUID) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM user_follows WHERE follower_id = $1 AND following_id = $2`,
		followerID, followingID)
	return err
}

func (r *followRepository) IsFollowing(ctx context.Context, followerID, followingID entities.UUID) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM user_follows WHERE follower_id = $1 AND following_id = $2`,
		followerID, followingID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *followRepository) GetFollowers(ctx context.Context, userID entities.UUID, limit, offset int) ([]*entities.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.avatar_url, u.role, u.is_verified, u.created_at, u.updated_at
		FROM users u
		INNER JOIN user_follows uf ON u.id = uf.follower_id
		WHERE uf.following_id = $1
		ORDER BY uf.followed_at DESC
		LIMIT $2 OFFSET $3`

	return r.scanUsers(ctx, query, userID, limit, offset)
}

func (r *followRepository) GetFollowing(ctx context.Context, userID entities.UUID, limit, offset int) ([]*entities.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.avatar_url, u.role, u.is_verified, u.created_at, u.updated_at
		FROM users u
		INNER JOIN user_follows uf ON u.id = uf.following_id
		WHERE uf.follower_id = $1
		ORDER BY uf.followed_at DESC
		LIMIT $2 OFFSET $3`

	return r.scanUsers(ctx, query, userID, limit, offset)
}

func (r *followRepository) GetFollowCount(ctx context.Context, userID entities.UUID) (followers int, following int, err error) {
	err = r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM user_follows WHERE following_id = $1`, userID).Scan(&followers)
	if err != nil {
		return 0, 0, err
	}

	err = r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM user_follows WHERE follower_id = $1`, userID).Scan(&following)
	if err != nil {
		return 0, 0, err
	}

	return followers, following, nil
}

func (r *followRepository) scanUsers(ctx context.Context, query string, args ...interface{}) ([]*entities.User, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
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
