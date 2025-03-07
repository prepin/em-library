package usecase

import (
	"context"
	"em-library/internal/entities"
)

type CreateSongUseCase interface {
	Execute(ctx context.Context, data entities.NewSongData) (*entities.SongData, error)
}

type createSongUseCase struct {
	transactionManager TransactionManager
	songRepo           SongRepo
	lyricsRepo         LyricsRepo
	songInfoService    SongInfoService
}

func NewCreateSongUseCase(
	tm TransactionManager,
	sr SongRepo,
	lr LyricsRepo,
	s SongInfoService,
) CreateSongUseCase {
	return &createSongUseCase{
		transactionManager: tm,
		songRepo:           sr,
		lyricsRepo:         lr,
		songInfoService:    s,
	}
}

func (u *createSongUseCase) Execute(ctx context.Context, data entities.NewSongData) (*entities.SongData, error) {
	info, err := u.songInfoService.GetInfo(ctx, data.Band, data.Song)
	if err != nil {
		return nil, err
	}

	data.ReleaseDate = info.ReleaseDate
	data.Link = info.Link

	var id int

	err = u.transactionManager.Do(ctx, func(ctx context.Context) error {
		id, err = u.songRepo.Create(ctx, data)
		if err != nil {
			return err
		}

		err = u.lyricsRepo.Create(ctx, entities.NewLyricsData{
			SongID:  id,
			Content: info.Lyrics,
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	song := entities.SongData{
		ID:          id,
		Band:        data.Band,
		Song:        data.Song,
		ReleaseDate: info.ReleaseDate,
		Link:        info.Link,
	}
	return &song, nil
}
