package repositories

import (
	"context"
	"nekosync/internal/domain/entities"
)

type FollowRepository interface {
	Follow(ctx context.Context, followerID, followingID entities.UUID) error
	Unfollow(ctx context.Context, followerID, followingID entities.UUID) error
	IsFollowing(ctx context.Context, followerID, followingID entities.UUID) (bool, error)
	GetFollowers(ctx context.Context, userID entities.UUID, limit, offset int) ([]*entities.User, error)
	GetFollowing(ctx context.Context, userID entities.UUID, limit, offset int) ([]*entities.User, error)
	GetFollowCount(ctx context.Context, userID entities.UUID) (followers int, following int, err error)
}
