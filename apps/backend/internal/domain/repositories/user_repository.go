package repositories

import (
	"context"
	"nekosync/internal/domain/entities"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id entities.UUID) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	GetByUsername(ctx context.Context, username string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id entities.UUID) error
	List(ctx context.Context, limit, offset int) ([]*entities.User, error)

	CreateProfile(ctx context.Context, profile *entities.UserProfile) error
	GetProfile(ctx context.Context, userID entities.UUID) (*entities.UserProfile, error)
	UpdateProfile(ctx context.Context, profile *entities.UserProfile) error
}
