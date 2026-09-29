package history

import (
	"context"
	"nekosync/internal/platform/entity"
)

type Repository interface {
	SaveWatchHistory(ctx context.Context, h *WatchHistory) error
	GetWatchHistory(ctx context.Context, userID entity.UUID, limit, offset int) ([]*WatchHistory, error)

	SaveReadHistory(ctx context.Context, h *ReadHistory) error
	GetReadHistory(ctx context.Context, userID entity.UUID, limit, offset int) ([]*ReadHistory, error)

	SaveListenHistory(ctx context.Context, h *ListenHistory) error
	GetListenHistory(ctx context.Context, userID entity.UUID, limit, offset int) ([]*ListenHistory, error)

	AddFavorite(ctx context.Context, f *Favorite) error
	RemoveFavorite(ctx context.Context, userID, contentID entity.UUID) error
	GetFavorites(ctx context.Context, userID entity.UUID, limit, offset int) ([]*Favorite, error)

	AddFavoriteMusic(ctx context.Context, f *FavoriteMusic) error
	RemoveFavoriteMusic(ctx context.Context, userID, musicID entity.UUID) error
	GetFavoriteMusic(ctx context.Context, userID entity.UUID, limit, offset int) ([]*FavoriteMusic, error)

	CreatePlaylist(ctx context.Context, p *Playlist) error
	GetPlaylistsByUserID(ctx context.Context, userID entity.UUID) ([]*Playlist, error)
	UpdatePlaylist(ctx context.Context, p *Playlist) error
	DeletePlaylist(ctx context.Context, id entity.UUID) error

	AddToPlaylist(ctx context.Context, pm *PlaylistMusic) error
	RemoveFromPlaylist(ctx context.Context, playlistID, musicID entity.UUID) error
	GetPlaylistMusic(ctx context.Context, playlistID entity.UUID) ([]*PlaylistMusic, error)
}
