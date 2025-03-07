package usecase

import (
	"context"
	"em-library/internal/entities"
)

type UpdateSongUseCase interface {
	Execute(ctx context.Context, songID int, data entities.UpdateSongData) error
}

type updateSongUseCase struct {
	transactionManager TransactionManager
	songRepo           SongRepo
	lyricsRepo         LyricsRepo
}

func NewUpdateSongUseCase(
	tm TransactionManager,
	sr SongRepo,
	lr LyricsRepo,
) UpdateSongUseCase {
	return &updateSongUseCase{
		transactionManager: tm,
		songRepo:           sr,
		lyricsRepo:         lr,
	}
}

func (u *updateSongUseCase) Execute(ctx context.Context, songID int, data entities.UpdateSongData) error {

	err := u.transactionManager.Do(ctx, func(ctx context.Context) error {
		if err := u.songRepo.Update(ctx, songID, data); err != nil {
			return err
		}

		if err := u.lyricsRepo.Update(ctx, songID, data); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
