package usecase

import (
	"context"
	"em-library/internal/entities"
)

type GetSongListUseCase interface {
	Execute(
		ctx context.Context,
		filter entities.SongFilterData,
	) ([]entities.SongData, error)
}

type getSongListUseCase struct {
	songRepo SongRepo
}

func NewGetSongListUseCase(sr SongRepo) GetSongListUseCase {
	return &getSongListUseCase{
		songRepo: sr,
	}
}

func (u *getSongListUseCase) Execute(
	ctx context.Context,
	filter entities.SongFilterData,
) ([]entities.SongData, error) {

	if filter.Limit == nil {
		limit := 50
		filter.Limit = &limit
	}

	if filter.Offset == nil {
		offset := 0
		filter.Offset = &offset
	}

	songData, err := u.songRepo.GetList(ctx, filter)
	if err != nil {
		return nil, err
	}

	return songData, nil
}
