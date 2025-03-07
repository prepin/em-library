package usecase

import (
	"context"
	"em-library/internal/entities"
	"strings"
)

type GetSongLyricsUseCase interface {
	Execute(
		ctx context.Context,
		songID int,
		filter entities.LyricsFilterData) ([]entities.LyricsVerseData, error)
}

type getSongLyricsUseCase struct {
	lyricsRepo LyricsRepo
}

func NewGetSongLyricsUsecase(lr LyricsRepo) GetSongLyricsUseCase {
	return &getSongLyricsUseCase{
		lyricsRepo: lr,
	}
}

func (u *getSongLyricsUseCase) Execute(
	ctx context.Context,
	songID int,
	filter entities.LyricsFilterData) ([]entities.LyricsVerseData, error) {

	lyrics, err := u.lyricsRepo.Get(ctx, songID)
	if err != nil {
		return nil, err
	}

	verses := strings.Split(lyrics.Content, "\\n\\n")

	result := make([]entities.LyricsVerseData, len(verses))

	for idx, verse := range verses {
		result[idx] = entities.LyricsVerseData{
			Index:   idx,
			Content: verse,
		}
	}

	var firstVerse int
	var lastVerse int

	if filter.Offset != nil {
		firstVerse = max(*filter.Offset, 0) // чтобы не взять отрицательный индекс в слайсе
		firstVerse = min(firstVerse, len(verses))
	} else {
		firstVerse = 0
	}

	if filter.Limit != nil {
		limit := max(0, *filter.Limit)
		lastVerse = min(firstVerse+limit, len(verses))
	} else {
		lastVerse = len(verses)
	}

	return result[firstVerse:lastVerse], nil
}
