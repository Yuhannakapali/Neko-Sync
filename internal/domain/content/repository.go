package content

import "nekosync/internal/domain"

type Repository interface {
	// Series operations
	CreateSeries(series *Series) error
	GetSeries(id domain.UUID) (*Series, error)
	UpdateSeries(series *Series) error
	DeleteSeries(id domain.UUID) error

	// Content operations
	CreateContent(content *Content) error
	GetContent(id domain.UUID) (*Content, error)
	GetContentsByType(contentType domain.ContentType) ([]*Content, error)
	GetContentsBySeries(seriesID domain.UUID) ([]*Content, error)
	UpdateContent(content *Content) error
	DeleteContent(id domain.UUID) error
	SearchContent(query string) ([]*Content, error)

	// Genre operations
	CreateGenre(genre *Genre) error
	GetGenre(id domain.UUID) (*Genre, error)
	GetAllGenres() ([]*Genre, error)
	UpdateGenre(genre *Genre) error
	DeleteGenre(id domain.UUID) error

	// Content Genre operations
	AddContentGenre(contentGenre *ContentGenre) error
	RemoveContentGenre(contentID, genreID domain.UUID) error
	GetContentGenres(contentID domain.UUID) ([]*Genre, error)

	// Tag operations
	CreateTag(tag *Tag) error
	GetTag(id domain.UUID) (*Tag, error)
	GetAllTags() ([]*Tag, error)
	UpdateTag(tag *Tag) error
	DeleteTag(id domain.UUID) error

	// Content Tag operations
	AddContentTag(contentTag *ContentTag) error
	RemoveContentTag(contentID, tagID domain.UUID) error
	GetContentTags(contentID domain.UUID) ([]*Tag, error)
}

type AnimeRepository interface {
	CreateAnime(anime *Anime) error
	GetAnime(contentID domain.UUID) (*Anime, error)
	UpdateAnime(anime *Anime) error
	DeleteAnime(contentID domain.UUID) error
	GetAnimeByStatus(status domain.ContentStatus) ([]*Anime, error)
	GetAnimeBySeasonAndYear(season domain.Season, year int) ([]*Anime, error)

	// Episode operations
	CreateEpisode(episode *Episode) error
	GetEpisode(id domain.UUID) (*Episode, error)
	GetEpisodesByContent(contentID domain.UUID) ([]*Episode, error)
	UpdateEpisode(episode *Episode) error
	DeleteEpisode(id domain.UUID) error
	IncrementEpisodeViews(id domain.UUID) error
}

type MangaRepository interface {
	CreateManga(manga *Manga) error
	GetManga(contentID domain.UUID) (*Manga, error)
	UpdateManga(manga *Manga) error
	DeleteManga(contentID domain.UUID) error
	GetMangaByStatus(status domain.ContentStatus) ([]*Manga, error)

	// Chapter operations
	CreateChapter(chapter *Chapter) error
	GetChapter(id domain.UUID) (*Chapter, error)
	GetChaptersByContent(contentID domain.UUID) ([]*Chapter, error)
	UpdateChapter(chapter *Chapter) error
	DeleteChapter(id domain.UUID) error
	IncrementChapterViews(id domain.UUID) error
}

type MovieRepository interface {
	CreateMovie(movie *Movie) error
	GetMovie(contentID domain.UUID) (*Movie, error)
	UpdateMovie(movie *Movie) error
	DeleteMovie(contentID domain.UUID) error
	GetMoviesByYear(year int) ([]*Movie, error)
	IncrementMovieViews(contentID domain.UUID) error
}

type MusicRepository interface {
	CreateMusic(music *Music) error
	GetMusic(contentID domain.UUID) (*Music, error)
	UpdateMusic(music *Music) error
	DeleteMusic(contentID domain.UUID) error
	GetMusicByType(musicType domain.MusicType) ([]*Music, error)
	IncrementMusicPlays(contentID domain.UUID) error

	// Artist operations
	CreateArtist(artist *Artist) error
	GetArtist(id domain.UUID) (*Artist, error)
	GetAllArtists() ([]*Artist, error)
	UpdateArtist(artist *Artist) error
	DeleteArtist(id domain.UUID) error
	SearchArtists(query string) ([]*Artist, error)

	// Music Artist operations
	AddMusicArtist(musicArtist *MusicArtist) error
	RemoveMusicArtist(musicID, artistID domain.UUID) error
	GetMusicArtists(musicID domain.UUID) ([]*Artist, error)
	GetArtistMusic(artistID domain.UUID) ([]*Music, error)
}
