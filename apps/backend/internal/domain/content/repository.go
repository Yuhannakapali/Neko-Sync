package content

import (
	"context"
	"nekosync/internal/domain/shared"
)

type Repository interface {
	Create(ctx context.Context, c *Content) error
	GetByID(ctx context.Context, id shared.UUID) (*Content, error)
	Update(ctx context.Context, c *Content) error
	Delete(ctx context.Context, id shared.UUID) error
	List(ctx context.Context, contentType Type, limit, offset int) ([]*Content, error)
	Search(ctx context.Context, query string, contentType Type, limit, offset int) ([]*Content, error)

	CreateSeries(ctx context.Context, s *Series) error
	GetSeriesByID(ctx context.Context, id shared.UUID) (*Series, error)
	GetContentBySeries(ctx context.Context, seriesID shared.UUID) ([]*Content, error)

	CreateGenre(ctx context.Context, g *Genre) error
	GetAllGenres(ctx context.Context) ([]*Genre, error)
	AddContentGenre(ctx context.Context, contentID, genreID shared.UUID) error
	RemoveContentGenre(ctx context.Context, contentID, genreID shared.UUID) error
	GetContentGenres(ctx context.Context, contentID shared.UUID) ([]*Genre, error)

	CreateTag(ctx context.Context, t *Tag) error
	GetAllTags(ctx context.Context) ([]*Tag, error)
	AddContentTag(ctx context.Context, contentID, tagID shared.UUID) error
	RemoveContentTag(ctx context.Context, contentID, tagID shared.UUID) error
	GetContentTags(ctx context.Context, contentID shared.UUID) ([]*Tag, error)
}

type AnimeRepository interface {
	CreateAnime(ctx context.Context, a *Anime) error
	GetAnimeByContentID(ctx context.Context, contentID shared.UUID) (*Anime, error)
	UpdateAnime(ctx context.Context, a *Anime) error
	GetAnimeByStatus(ctx context.Context, status Status, limit, offset int) ([]*Anime, error)

	CreateEpisode(ctx context.Context, e *Episode) error
	GetEpisodesByContentID(ctx context.Context, contentID shared.UUID) ([]*Episode, error)
	GetEpisodeByNumber(ctx context.Context, contentID shared.UUID, episodeNumber int) (*Episode, error)
	UpdateEpisode(ctx context.Context, e *Episode) error
	DeleteEpisode(ctx context.Context, id shared.UUID) error
}

type MangaRepository interface {
	CreateManga(ctx context.Context, m *Manga) error
	GetMangaByContentID(ctx context.Context, contentID shared.UUID) (*Manga, error)
	UpdateManga(ctx context.Context, m *Manga) error

	CreateChapter(ctx context.Context, c *Chapter) error
	GetChaptersByContentID(ctx context.Context, contentID shared.UUID) ([]*Chapter, error)
	GetChapterByNumber(ctx context.Context, contentID shared.UUID, chapterNumber int) (*Chapter, error)
	UpdateChapter(ctx context.Context, c *Chapter) error
	DeleteChapter(ctx context.Context, id shared.UUID) error
}

type MovieRepository interface {
	CreateMovie(ctx context.Context, m *Movie) error
	GetMovieByContentID(ctx context.Context, contentID shared.UUID) (*Movie, error)
	UpdateMovie(ctx context.Context, m *Movie) error
}

type MusicRepository interface {
	CreateMusic(ctx context.Context, m *Music) error
	GetMusicByContentID(ctx context.Context, contentID shared.UUID) (*Music, error)
	UpdateMusic(ctx context.Context, m *Music) error
	GetMusicByType(ctx context.Context, musicType MusicType, limit, offset int) ([]*Music, error)

	CreateArtist(ctx context.Context, a *Artist) error
	GetArtistByID(ctx context.Context, id shared.UUID) (*Artist, error)
	UpdateArtist(ctx context.Context, a *Artist) error
	AddMusicArtist(ctx context.Context, musicID, artistID shared.UUID, role *string) error
	GetArtistsByMusicID(ctx context.Context, musicID shared.UUID) ([]*Artist, error)
	GetMusicByArtistID(ctx context.Context, artistID shared.UUID) ([]*Music, error)
}
