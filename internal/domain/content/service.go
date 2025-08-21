package content

import (
	"errors"
	"nekosync/internal/domain"
)

type Service struct {
	repo      Repository
	animeRepo AnimeRepository
	mangaRepo MangaRepository
	movieRepo MovieRepository
	musicRepo MusicRepository
}

func NewService(repo Repository, animeRepo AnimeRepository, mangaRepo MangaRepository, movieRepo MovieRepository, musicRepo MusicRepository) *Service {
	return &Service{
		repo:      repo,
		animeRepo: animeRepo,
		mangaRepo: mangaRepo,
		movieRepo: movieRepo,
		musicRepo: musicRepo,
	}
}

// Content operations
func (s *Service) CreateContent(content *Content) error {
	if content.Title == "" {
		return errors.New("title is required")
	}
	return s.repo.CreateContent(content)
}

func (s *Service) GetContent(id domain.UUID) (*Content, error) {
	return s.repo.GetContent(id)
}

func (s *Service) SearchContent(query string) ([]*Content, error) {
	if query == "" {
		return nil, errors.New("search query cannot be empty")
	}
	return s.repo.SearchContent(query)
}

func (s *Service) GetContentsByType(contentType domain.ContentType) ([]*Content, error) {
	return s.repo.GetContentsByType(contentType)
}

// Series operations
func (s *Service) CreateSeries(series *Series) error {
	if series.Title == "" {
		return errors.New("series title is required")
	}
	return s.repo.CreateSeries(series)
}

func (s *Service) GetSeries(id domain.UUID) (*Series, error) {
	return s.repo.GetSeries(id)
}

// Genre operations
func (s *Service) CreateGenre(genre *Genre) error {
	if genre.Name == "" {
		return errors.New("genre name is required")
	}
	return s.repo.CreateGenre(genre)
}

func (s *Service) GetAllGenres() ([]*Genre, error) {
	return s.repo.GetAllGenres()
}

func (s *Service) AddContentGenre(contentID, genreID domain.UUID) error {
	contentGenre := &ContentGenre{
		ContentID: contentID,
		GenreID:   genreID,
	}
	return s.repo.AddContentGenre(contentGenre)
}

// Anime operations
func (s *Service) CreateAnime(anime *Anime) error {
	return s.animeRepo.CreateAnime(anime)
}

func (s *Service) GetAnime(contentID domain.UUID) (*Anime, error) {
	return s.animeRepo.GetAnime(contentID)
}

func (s *Service) CreateEpisode(episode *Episode) error {
	if episode.Title == "" {
		return errors.New("episode title is required")
	}
	return s.animeRepo.CreateEpisode(episode)
}

func (s *Service) GetEpisodesByContent(contentID domain.UUID) ([]*Episode, error) {
	return s.animeRepo.GetEpisodesByContent(contentID)
}

func (s *Service) WatchEpisode(episodeID domain.UUID) error {
	return s.animeRepo.IncrementEpisodeViews(episodeID)
}

// Manga operations
func (s *Service) CreateManga(manga *Manga) error {
	return s.mangaRepo.CreateManga(manga)
}

func (s *Service) GetManga(contentID domain.UUID) (*Manga, error) {
	return s.mangaRepo.GetManga(contentID)
}

func (s *Service) CreateChapter(chapter *Chapter) error {
	if chapter.Title == "" {
		return errors.New("chapter title is required")
	}
	return s.mangaRepo.CreateChapter(chapter)
}

func (s *Service) GetChaptersByContent(contentID domain.UUID) ([]*Chapter, error) {
	return s.mangaRepo.GetChaptersByContent(contentID)
}

func (s *Service) ReadChapter(chapterID domain.UUID) error {
	return s.mangaRepo.IncrementChapterViews(chapterID)
}

// Movie operations
func (s *Service) CreateMovie(movie *Movie) error {
	return s.movieRepo.CreateMovie(movie)
}

func (s *Service) GetMovie(contentID domain.UUID) (*Movie, error) {
	return s.movieRepo.GetMovie(contentID)
}

func (s *Service) WatchMovie(contentID domain.UUID) error {
	return s.movieRepo.IncrementMovieViews(contentID)
}

// Music operations
func (s *Service) CreateMusic(music *Music) error {
	return s.musicRepo.CreateMusic(music)
}

func (s *Service) GetMusic(contentID domain.UUID) (*Music, error) {
	return s.musicRepo.GetMusic(contentID)
}

func (s *Service) PlayMusic(contentID domain.UUID) error {
	return s.musicRepo.IncrementMusicPlays(contentID)
}

func (s *Service) CreateArtist(artist *Artist) error {
	if artist.Name == "" {
		return errors.New("artist name is required")
	}
	return s.musicRepo.CreateArtist(artist)
}

func (s *Service) GetArtist(id domain.UUID) (*Artist, error) {
	return s.musicRepo.GetArtist(id)
}

func (s *Service) SearchArtists(query string) ([]*Artist, error) {
	if query == "" {
		return nil, errors.New("search query cannot be empty")
	}
	return s.musicRepo.SearchArtists(query)
}
