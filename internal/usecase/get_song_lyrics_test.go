package usecase_test

import (
	"context"
	"em-library/internal/entities"
	"em-library/internal/usecase"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSongLyricsUseCase_Execute_Success(t *testing.T) {
	mockLyricsRepo := new(MockLyricsRepo)
	useCase := usecase.NewGetSongLyricsUsecase(mockLyricsRepo)

	ctx := context.Background()
	songID := 123

	mockLyrics := entities.LyricsData{
		SongID:  songID,
		Content: "First verse line 1\\nFirst verse line 2\\n\\nSecond verse line 1\\nSecond verse line 2\\n\\nThird verse",
	}

	mockLyricsRepo.On("Get", ctx, songID).Return(mockLyrics, nil)

	result, err := useCase.Execute(ctx, songID, entities.LyricsFilterData{})

	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, 0, result[0].Index)
	assert.Equal(t, "First verse line 1\\nFirst verse line 2", result[0].Content)
	assert.Equal(t, 1, result[1].Index)
	assert.Equal(t, "Second verse line 1\\nSecond verse line 2", result[1].Content)
	assert.Equal(t, 2, result[2].Index)
	assert.Equal(t, "Third verse", result[2].Content)

	mockLyricsRepo.AssertExpectations(t)
}

func TestGetSongLyricsUseCase_Execute_WithFilter(t *testing.T) {
	mockLyricsRepo := new(MockLyricsRepo)
	useCase := usecase.NewGetSongLyricsUsecase(mockLyricsRepo)

	ctx := context.Background()
	songID := 123

	mockLyrics := entities.LyricsData{
		SongID:  songID,
		Content: "Verse 1\\n\\nVerse 2\\n\\nVerse 3\\n\\nVerse 4\\n\\nVerse 5",
	}

	mockLyricsRepo.On("Get", ctx, songID).Return(mockLyrics, nil)

	offset := 1
	limit := 2
	filter := entities.LyricsFilterData{
		Offset: &offset,
		Limit:  &limit,
	}

	result, err := useCase.Execute(ctx, songID, filter)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 1, result[0].Index)
	assert.Equal(t, "Verse 2", result[0].Content)
	assert.Equal(t, 2, result[1].Index)
	assert.Equal(t, "Verse 3", result[1].Content)

	mockLyricsRepo.AssertExpectations(t)
}

func TestGetSongLyricsUseCase_Execute_OffsetBeyondLength(t *testing.T) {
	mockLyricsRepo := new(MockLyricsRepo)
	useCase := usecase.NewGetSongLyricsUsecase(mockLyricsRepo)

	ctx := context.Background()
	songID := 123

	mockLyrics := entities.LyricsData{
		SongID:  songID,
		Content: "Verse 1\\n\\nVerse 2\\n\\nVerse 3",
	}

	mockLyricsRepo.On("Get", ctx, songID).Return(mockLyrics, nil)

	offset := 10
	filter := entities.LyricsFilterData{
		Offset: &offset,
	}

	result, err := useCase.Execute(ctx, songID, filter)

	assert.NoError(t, err)
	assert.Len(t, result, 0)

	mockLyricsRepo.AssertExpectations(t)
}

func TestGetSongLyricsUseCase_Execute_LimitBeyondLength(t *testing.T) {
	mockLyricsRepo := new(MockLyricsRepo)
	useCase := usecase.NewGetSongLyricsUsecase(mockLyricsRepo)

	ctx := context.Background()
	songID := 123

	mockLyrics := entities.LyricsData{
		SongID:  songID,
		Content: "Verse 1\\n\\nVerse 2\\n\\nVerse 3",
	}

	mockLyricsRepo.On("Get", ctx, songID).Return(mockLyrics, nil)

	offset := 1
	limit := 10
	filter := entities.LyricsFilterData{
		Offset: &offset,
		Limit:  &limit,
	}

	result, err := useCase.Execute(ctx, songID, filter)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Verse 2", result[0].Content)
	assert.Equal(t, "Verse 3", result[1].Content)

	mockLyricsRepo.AssertExpectations(t)
}

func TestGetSongLyricsUseCase_Execute_RepoError(t *testing.T) {
	mockLyricsRepo := new(MockLyricsRepo)
	useCase := usecase.NewGetSongLyricsUsecase(mockLyricsRepo)

	ctx := context.Background()
	songID := 123

	expectedError := errors.New("repository error")
	mockLyricsRepo.On("Get", ctx, songID).Return(entities.LyricsData{}, expectedError)

	result, err := useCase.Execute(ctx, songID, entities.LyricsFilterData{})

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)

	mockLyricsRepo.AssertExpectations(t)
}

func TestGetSongLyricsUseCase_Execute_EmptyContent(t *testing.T) {
	mockLyricsRepo := new(MockLyricsRepo)
	useCase := usecase.NewGetSongLyricsUsecase(mockLyricsRepo)

	ctx := context.Background()
	songID := 123

	mockLyrics := entities.LyricsData{
		SongID:  songID,
		Content: "",
	}

	mockLyricsRepo.On("Get", ctx, songID).Return(mockLyrics, nil)

	result, err := useCase.Execute(ctx, songID, entities.LyricsFilterData{})

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "", result[0].Content)

	mockLyricsRepo.AssertExpectations(t)
}

func TestGetSongLyricsUseCase_Execute_NilOffsetAndLimit(t *testing.T) {
	mockLyricsRepo := new(MockLyricsRepo)
	useCase := usecase.NewGetSongLyricsUsecase(mockLyricsRepo)

	ctx := context.Background()
	songID := 123

	mockLyrics := entities.LyricsData{
		SongID:  songID,
		Content: "Verse 1\\n\\nVerse 2\\n\\nVerse 3",
	}

	mockLyricsRepo.On("Get", ctx, songID).Return(mockLyrics, nil)

	filter := entities.LyricsFilterData{
		Offset: nil,
		Limit:  nil,
	}

	result, err := useCase.Execute(ctx, songID, filter)

	assert.NoError(t, err)
	assert.Len(t, result, 3)

	mockLyricsRepo.AssertExpectations(t)
}

func TestGetSongLyricsUseCase_Execute_LimitOnly(t *testing.T) {
	mockLyricsRepo := new(MockLyricsRepo)
	useCase := usecase.NewGetSongLyricsUsecase(mockLyricsRepo)

	ctx := context.Background()
	songID := 123

	mockLyrics := entities.LyricsData{
		SongID:  songID,
		Content: "Verse 1\\n\\nVerse 2\\n\\nVerse 3\\n\\nVerse 4",
	}

	mockLyricsRepo.On("Get", ctx, songID).Return(mockLyrics, nil)

	limit := 2
	filter := entities.LyricsFilterData{
		Limit: &limit,
	}

	result, err := useCase.Execute(ctx, songID, filter)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Verse 1", result[0].Content)
	assert.Equal(t, "Verse 2", result[1].Content)

	mockLyricsRepo.AssertExpectations(t)
}
