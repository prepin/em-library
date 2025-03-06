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
}

type LyricsRepo interface {
	Create(ctx context.Context, data entities.NewLyricsData) error
}

type SongInfoService interface {
	GetInfo(ctx context.Context, group, song string) (*entities.SongDetail, error)
}
