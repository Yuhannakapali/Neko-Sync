package repositories

import (
	"context"
	"nekosync/internal/domain/entities"
)

type MusicRepository interface {
	CreateMusic(ctx context.Context, music *entities.Music) error
	GetMusicByContentID(ctx context.Context, contentID entities.UUID) (*entities.Music, error)
	UpdateMusic(ctx context.Context, music *entities.Music) error
	GetMusicByType(ctx context.Context, musicType entities.MusicType, limit, offset int) ([]*entities.Music, error)

	CreateArtist(ctx context.Context, artist *entities.Artist) error
	GetArtistByID(ctx context.Context, id entities.UUID) (*entities.Artist, error)
	UpdateArtist(ctx context.Context, artist *entities.Artist) error
	AddMusicArtist(ctx context.Context, musicID, artistID entities.UUID, role *string) error
	GetArtistsByMusicID(ctx context.Context, musicID entities.UUID) ([]*entities.Artist, error)
	GetMusicByArtistID(ctx context.Context, artistID entities.UUID) ([]*entities.Music, error)
}
