package repositories

import (
	"context"
	"nekosync/internal/domain/entities"
)

// ContentRepository defines the contract for content data operations
type ContentRepository interface {
	// Content CRUD operations
	Create(ctx context.Context, content *entities.Content) error
	GetByID(ctx context.Context, id entities.UUID) (*entities.Content, error)
	Update(ctx context.Context, content *entities.Content) error
	Delete(ctx context.Context, id entities.UUID) error
	List(ctx context.Context, contentType entities.ContentType, limit, offset int) ([]*entities.Content, error)
	Search(ctx context.Context, query string, contentType entities.ContentType, limit, offset int) ([]*entities.Content, error)

	// Series operations
	CreateSeries(ctx context.Context, series *entities.Series) error
	GetSeriesByID(ctx context.Context, id entities.UUID) (*entities.Series, error)
	GetContentBySeries(ctx context.Context, seriesID entities.UUID) ([]*entities.Content, error)

	// Genre operations
	CreateGenre(ctx context.Context, genre *entities.Genre) error
	GetAllGenres(ctx context.Context) ([]*entities.Genre, error)
	AddContentGenre(ctx context.Context, contentID, genreID entities.UUID) error
	RemoveContentGenre(ctx context.Context, contentID, genreID entities.UUID) error
	GetContentGenres(ctx context.Context, contentID entities.UUID) ([]*entities.Genre, error)

	// Tag operations
	CreateTag(ctx context.Context, tag *entities.Tag) error
	GetAllTags(ctx context.Context) ([]*entities.Tag, error)
	AddContentTag(ctx context.Context, contentID, tagID entities.UUID) error
	RemoveContentTag(ctx context.Context, contentID, tagID entities.UUID) error
	GetContentTags(ctx context.Context, contentID entities.UUID) ([]*entities.Tag, error)

	// Anime operations
	CreateAnime(ctx context.Context, anime *entities.Anime) error
	GetAnimeByContentID(ctx context.Context, contentID entities.UUID) (*entities.Anime, error)
	UpdateAnime(ctx context.Context, anime *entities.Anime) error
	GetAnimeByStatus(ctx context.Context, status entities.ContentStatus, limit, offset int) ([]*entities.Anime, error)

	// Episode operations
	CreateEpisode(ctx context.Context, episode *entities.Episode) error
	GetEpisodesByContentID(ctx context.Context, contentID entities.UUID) ([]*entities.Episode, error)
	GetEpisodeByNumber(ctx context.Context, contentID entities.UUID, episodeNumber int) (*entities.Episode, error)
	UpdateEpisode(ctx context.Context, episode *entities.Episode) error
	DeleteEpisode(ctx context.Context, id entities.UUID) error

	// Manga operations
	CreateManga(ctx context.Context, manga *entities.Manga) error
	GetMangaByContentID(ctx context.Context, contentID entities.UUID) (*entities.Manga, error)
	UpdateManga(ctx context.Context, manga *entities.Manga) error

	// Chapter operations
	CreateChapter(ctx context.Context, chapter *entities.Chapter) error
	GetChaptersByContentID(ctx context.Context, contentID entities.UUID) ([]*entities.Chapter, error)
	GetChapterByNumber(ctx context.Context, contentID entities.UUID, chapterNumber int) (*entities.Chapter, error)
	UpdateChapter(ctx context.Context, chapter *entities.Chapter) error
	DeleteChapter(ctx context.Context, id entities.UUID) error

	// Movie operations
	CreateMovie(ctx context.Context, movie *entities.Movie) error
	GetMovieByContentID(ctx context.Context, contentID entities.UUID) (*entities.Movie, error)
	UpdateMovie(ctx context.Context, movie *entities.Movie) error

	// Music operations
	CreateMusic(ctx context.Context, music *entities.Music) error
	GetMusicByContentID(ctx context.Context, contentID entities.UUID) (*entities.Music, error)
	UpdateMusic(ctx context.Context, music *entities.Music) error
	GetMusicByType(ctx context.Context, musicType entities.MusicType, limit, offset int) ([]*entities.Music, error)

	// Artist operations
	CreateArtist(ctx context.Context, artist *entities.Artist) error
	GetArtistByID(ctx context.Context, id entities.UUID) (*entities.Artist, error)
	UpdateArtist(ctx context.Context, artist *entities.Artist) error
	AddMusicArtist(ctx context.Context, musicID, artistID entities.UUID, role *string) error
	GetArtistsByMusicID(ctx context.Context, musicID entities.UUID) ([]*entities.Artist, error)
	GetMusicByArtistID(ctx context.Context, artistID entities.UUID) ([]*entities.Music, error)
}
