package repositories

import (
	"context"
	"nekosync/internal/domain/entities"
)

type MovieRepository interface {
	CreateMovie(ctx context.Context, movie *entities.Movie) error
	GetMovieByContentID(ctx context.Context, contentID entities.UUID) (*entities.Movie, error)
	UpdateMovie(ctx context.Context, movie *entities.Movie) error
}
