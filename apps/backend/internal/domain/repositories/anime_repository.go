package repositories

import (
	"context"
	"nekosync/internal/domain/entities"
)

type AnimeRepository interface {
	CreateAnime(ctx context.Context, anime *entities.Anime) error
	GetAnimeByContentID(ctx context.Context, contentID entities.UUID) (*entities.Anime, error)
	UpdateAnime(ctx context.Context, anime *entities.Anime) error
	GetAnimeByStatus(ctx context.Context, status entities.ContentStatus, limit, offset int) ([]*entities.Anime, error)

	CreateEpisode(ctx context.Context, episode *entities.Episode) error
	GetEpisodesByContentID(ctx context.Context, contentID entities.UUID) ([]*entities.Episode, error)
	GetEpisodeByNumber(ctx context.Context, contentID entities.UUID, episodeNumber int) (*entities.Episode, error)
	UpdateEpisode(ctx context.Context, episode *entities.Episode) error
	DeleteEpisode(ctx context.Context, id entities.UUID) error
}
