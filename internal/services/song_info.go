package services

import (
	"context"
	"em-library/internal/entities"
	"time"
)

type RESTSongInfoService struct {
	// TODO: add logger
}

func NewRESTSongInfoService() *RESTSongInfoService {
	return &RESTSongInfoService{}
}

func (s *RESTSongInfoService) GetInfo(ctx context.Context, group, song string) (*entities.SongDetail, error) {

	// TODO: логи всех уровней

	songDetail := entities.SongDetail{
		ReleaseDate: time.Now(),
		Text: `Ooh baby, don't you know I suffer?\nOoh baby, can
you hear me moan?\nYou caught me under false pretenses\nHow long
before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set
my soul alight`,
		Link: "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}
	return &songDetail, nil
}
