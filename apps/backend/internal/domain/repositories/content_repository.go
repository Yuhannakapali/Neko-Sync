package repositories

import (
	"context"
	"nekosync/internal/domain/entities"
)

type ContentRepository interface {
	Create(ctx context.Context, content *entities.Content) error
	GetByID(ctx context.Context, id entities.UUID) (*entities.Content, error)
	Update(ctx context.Context, content *entities.Content) error
	Delete(ctx context.Context, id entities.UUID) error
	List(ctx context.Context, contentType entities.ContentType, limit, offset int) ([]*entities.Content, error)
	Search(ctx context.Context, query string, contentType entities.ContentType, limit, offset int) ([]*entities.Content, error)

	CreateSeries(ctx context.Context, series *entities.Series) error
	GetSeriesByID(ctx context.Context, id entities.UUID) (*entities.Series, error)
	GetContentBySeries(ctx context.Context, seriesID entities.UUID) ([]*entities.Content, error)

	CreateGenre(ctx context.Context, genre *entities.Genre) error
	GetAllGenres(ctx context.Context) ([]*entities.Genre, error)
	AddContentGenre(ctx context.Context, contentID, genreID entities.UUID) error
	RemoveContentGenre(ctx context.Context, contentID, genreID entities.UUID) error
	GetContentGenres(ctx context.Context, contentID entities.UUID) ([]*entities.Genre, error)

	CreateTag(ctx context.Context, tag *entities.Tag) error
	GetAllTags(ctx context.Context) ([]*entities.Tag, error)
	AddContentTag(ctx context.Context, contentID, tagID entities.UUID) error
	RemoveContentTag(ctx context.Context, contentID, tagID entities.UUID) error
	GetContentTags(ctx context.Context, contentID entities.UUID) ([]*entities.Tag, error)
}
