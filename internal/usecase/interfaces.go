package usecase

import (
	"context"
	"em-library/internal/entities"
)

type TransactionManager interface {
	Do(ctx context.Context, f func(ctx context.Context) error) error
}

type Repos struct {
	TransactionManager TransactionManager
	SongRepo           SongRepo
	LyricsRepo         LyricsRepo
}

type Services struct {
	SongInfoService SongInfoService
}

type SongRepo interface {
	Create(ctx context.Context, data entities.NewSongData) (int, error)
	GetList(ctx context.Context, filter entities.SongFilterData) ([]entities.SongData, error)
	Delete(ctx context.Context, songID int) error
}

type LyricsRepo interface {
	Create(ctx context.Context, data entities.NewLyricsData) error
	Get(ctx context.Context, songID int) (entities.LyricsData, error)
	Delete(ctx context.Context, songID int) error
}

type SongInfoService interface {
	GetInfo(ctx context.Context, group, song string) (*entities.SongDetail, error)
}
