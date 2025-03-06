package usecase

import "context"

type DeleteSongUseCase interface {
	Execute(ctx context.Context, songID int) error
}

type deleteSongUseCase struct {
	transactionManager TransactionManager
	songRepo           SongRepo
	lyricsRepo         LyricsRepo
}

func NewDeleteSongUseCase(tm TransactionManager, sr SongRepo, lr LyricsRepo) DeleteSongUseCase {
	return &deleteSongUseCase{
		transactionManager: tm,
		songRepo:           sr,
		lyricsRepo:         lr,
	}
}

func (u *deleteSongUseCase) Execute(ctx context.Context, songID int) error {

	err := u.transactionManager.Do(ctx, func(ctx context.Context) error {

		if err := u.lyricsRepo.Delete(ctx, songID); err != nil {
			return err
		}

		if err := u.songRepo.Delete(ctx, songID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
