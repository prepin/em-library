package usecase

import (
	"context"
	"em-library/internal/entities"
)

type Repos struct {
	SongRepo SongRepo
}

type Services struct {
	SongInfoService SongInfoService
}

type SongRepo interface {
	Create(ctx context.Context, data entities.NewSongData) (int, error)
}

type SongInfoService interface {
	GetInfo(ctx context.Context, group, song string) (*entities.SongDetail, error)
}
