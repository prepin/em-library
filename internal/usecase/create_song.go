package usecase

import (
	"context"
	"em-library/internal/entities"
)

type CreateSongUseCase interface {
	Execute(ctx context.Context, data entities.NewSongData) (*entities.SongData, error)
}

type createSongUseCase struct {
	songRepo        SongRepo
	songInfoService SongInfoService
}

func NewCreateSongUseCase(r SongRepo, s SongInfoService) CreateSongUseCase {
	return &createSongUseCase{
		songRepo:        r,
		songInfoService: s,
	}
}

func (u *createSongUseCase) Execute(ctx context.Context, data entities.NewSongData) (*entities.SongData, error) {
	// TODO: сходить в сторонний сервис за дополнительными данными
	// не забыть передать внутрь контекст и внутри сервиса отменять по таймауту
	info, err := u.songInfoService.GetInfo(ctx, data.Group, data.Song)
	if err != nil {
		return nil, err
	}

	id, err := u.songRepo.Create(ctx, data)
	if err != nil {
		return nil, err
	}

	// TODO: добавить данные из стороннего сервиса в итоговую структуру
	song := entities.SongData{
		ID:          id,
		Group:       data.Group,
		Song:        data.Song,
		ReleaseDate: info.ReleaseDate,
		Link:        info.Link,
	}
	return &song, nil
}
