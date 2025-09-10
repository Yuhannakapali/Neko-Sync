package history

import (
	"errors"

	"nekosync/internal/domain"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Watch History operations
func (s *Service) UpdateWatchProgress(userID, contentID, episodeID domain.UUID, progressSeconds int, finished bool) error {
	if progressSeconds < 0 {
		return errors.New("progress cannot be negative")
	}

	history := &WatchHistory{
		UserID:          userID,
		ContentID:       contentID,
		EpisodeID:       episodeID,
		ProgressSeconds: progressSeconds,
		Finished:        finished,
		Status:          "watching",
	}

	// Check if history exists
	existing, _ := s.repo.GetWatchHistory(userID, contentID, episodeID)
	if existing != nil {
		return s.repo.UpdateWatchHistory(history)
	}

	return s.repo.CreateWatchHistory(history)
}

func (s *Service) GetUserWatchHistory(userID domain.UUID, limit, offset int) ([]*WatchHistory, error) {
	if limit <= 0 {
		limit = 50 // default limit
	}
	return s.repo.GetUserWatchHistory(userID, limit, offset)
}

// Read History operations
func (s *Service) UpdateReadProgress(userID, contentID, chapterID domain.UUID, pageNumber int, finished bool) error {
	if pageNumber < 0 {
		return errors.New("page number cannot be negative")
	}

	history := &ReadHistory{
		UserID:     userID,
		ContentID:  contentID,
		ChapterID:  chapterID,
		PageNumber: pageNumber,
		Finished:   finished,
		Status:     "reading",
	}

	// Check if history exists
	existing, _ := s.repo.GetReadHistory(userID, contentID, chapterID)
	if existing != nil {
		return s.repo.UpdateReadHistory(history)
	}

	return s.repo.CreateReadHistory(history)
}

func (s *Service) GetUserReadHistory(userID domain.UUID, limit, offset int) ([]*ReadHistory, error) {
	if limit <= 0 {
		limit = 50 // default limit
	}
	return s.repo.GetUserReadHistory(userID, limit, offset)
}

// Listen History operations
func (s *Service) UpdateListenProgress(userID, musicID domain.UUID, progressSeconds int, finished bool) error {
	if progressSeconds < 0 {
		return errors.New("progress cannot be negative")
	}

	history := &ListenHistory{
		UserID:          userID,
		MusicID:         musicID,
		ProgressSeconds: progressSeconds,
		Finished:        finished,
		Status:          "listening",
	}

	// Check if history exists
	existing, _ := s.repo.GetListenHistory(userID, musicID)
	if existing != nil {
		return s.repo.UpdateListenHistory(history)
	}

	return s.repo.CreateListenHistory(history)
}

func (s *Service) GetUserListenHistory(userID domain.UUID, limit, offset int) ([]*ListenHistory, error) {
	if limit <= 0 {
		limit = 50 // default limit
	}
	return s.repo.GetUserListenHistory(userID, limit, offset)
}

// Favorites operations
func (s *Service) AddToFavorites(userID, contentID domain.UUID) error {
	// Check if already favorited
	isFav, err := s.repo.IsFavorite(userID, contentID)
	if err != nil {
		return err
	}
	if isFav {
		return errors.New("content already in favorites")
	}

	favorite := &Favorite{
		UserID:    userID,
		ContentID: contentID,
	}

	return s.repo.CreateFavorite(favorite)
}

func (s *Service) RemoveFromFavorites(userID, contentID domain.UUID) error {
	return s.repo.DeleteFavorite(userID, contentID)
}

func (s *Service) GetUserFavorites(userID domain.UUID) ([]*Favorite, error) {
	return s.repo.GetUserFavorites(userID)
}

// Favorite Music operations
func (s *Service) AddMusicToFavorites(userID, musicID domain.UUID) error {
	// Check if already favorited
	isFav, err := s.repo.IsFavoriteMusic(userID, musicID)
	if err != nil {
		return err
	}
	if isFav {
		return errors.New("music already in favorites")
	}

	favorite := &FavoriteMusic{
		UserID:  userID,
		MusicID: musicID,
	}

	return s.repo.CreateFavoriteMusic(favorite)
}

func (s *Service) RemoveMusicFromFavorites(userID, musicID domain.UUID) error {
	return s.repo.DeleteFavoriteMusic(userID, musicID)
}

// Playlist operations
func (s *Service) CreatePlaylist(playlist *Playlist) error {
	if playlist.Name == "" {
		return errors.New("playlist name is required")
	}
	return s.repo.CreatePlaylist(playlist)
}

func (s *Service) GetPlaylist(id domain.UUID) (*Playlist, error) {
	return s.repo.GetPlaylist(id)
}

func (s *Service) GetUserPlaylists(userID domain.UUID) ([]*Playlist, error) {
	return s.repo.GetUserPlaylists(userID)
}

func (s *Service) AddMusicToPlaylist(playlistID, musicID domain.UUID) error {
	// Get current playlist music to determine order
	playlistMusic, err := s.repo.GetPlaylistMusic(playlistID)
	if err != nil {
		return err
	}

	newEntry := &PlaylistMusic{
		PlaylistID: playlistID,
		MusicID:    musicID,
		OrderIndex: len(playlistMusic) + 1,
	}

	return s.repo.AddPlaylistMusic(newEntry)
}

func (s *Service) RemoveMusicFromPlaylist(playlistID, musicID domain.UUID) error {
	return s.repo.RemovePlaylistMusic(playlistID, musicID)
}

func (s *Service) GetPlaylistMusic(playlistID domain.UUID) ([]*PlaylistMusic, error) {
	return s.repo.GetPlaylistMusic(playlistID)
}
