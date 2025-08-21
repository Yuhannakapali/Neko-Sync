package history

import "nekosync/internal/domain"

type Repository interface {
	// Watch History operations
	CreateWatchHistory(history *WatchHistory) error
	GetWatchHistory(userID, contentID, episodeID domain.UUID) (*WatchHistory, error)
	GetUserWatchHistory(userID domain.UUID, limit, offset int) ([]*WatchHistory, error)
	UpdateWatchHistory(history *WatchHistory) error
	DeleteWatchHistory(userID, contentID, episodeID domain.UUID) error

	// Read History operations
	CreateReadHistory(history *ReadHistory) error
	GetReadHistory(userID, contentID, chapterID domain.UUID) (*ReadHistory, error)
	GetUserReadHistory(userID domain.UUID, limit, offset int) ([]*ReadHistory, error)
	UpdateReadHistory(history *ReadHistory) error
	DeleteReadHistory(userID, contentID, chapterID domain.UUID) error

	// Listen History operations
	CreateListenHistory(history *ListenHistory) error
	GetListenHistory(userID, musicID domain.UUID) (*ListenHistory, error)
	GetUserListenHistory(userID domain.UUID, limit, offset int) ([]*ListenHistory, error)
	UpdateListenHistory(history *ListenHistory) error
	DeleteListenHistory(userID, musicID domain.UUID) error

	// Favorites operations
	CreateFavorite(favorite *Favorite) error
	DeleteFavorite(userID, contentID domain.UUID) error
	GetUserFavorites(userID domain.UUID) ([]*Favorite, error)
	IsFavorite(userID, contentID domain.UUID) (bool, error)

	// Favorite Music operations
	CreateFavoriteMusic(favorite *FavoriteMusic) error
	DeleteFavoriteMusic(userID, musicID domain.UUID) error
	GetUserFavoriteMusic(userID domain.UUID) ([]*FavoriteMusic, error)
	IsFavoriteMusic(userID, musicID domain.UUID) (bool, error)

	// Playlist operations
	CreatePlaylist(playlist *Playlist) error
	GetPlaylist(id domain.UUID) (*Playlist, error)
	GetUserPlaylists(userID domain.UUID) ([]*Playlist, error)
	GetPublicPlaylists() ([]*Playlist, error)
	UpdatePlaylist(playlist *Playlist) error
	DeletePlaylist(id domain.UUID) error

	// Playlist Music operations
	AddPlaylistMusic(playlistMusic *PlaylistMusic) error
	RemovePlaylistMusic(playlistID, musicID domain.UUID) error
	GetPlaylistMusic(playlistID domain.UUID) ([]*PlaylistMusic, error)
	ReorderPlaylistMusic(playlistID domain.UUID, musicIDs []domain.UUID) error
}
