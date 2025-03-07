package usecase_test

import (
	"context"
	"em-library/internal/entities"
	"em-library/internal/usecase"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateSongUseCase_Execute_Success(t *testing.T) {
	mockTM := new(MockTransactionManager)
	mockSongRepo := new(MockSongRepo)
	mockLyricsRepo := new(MockLyricsRepo)
	useCase := usecase.NewUpdateSongUseCase(mockTM, mockSongRepo, mockLyricsRepo)

	ctx := context.Background()
	songID := 123
	band := "Updated Band"
	song := "Updated Song"
	link := "https://updated.example.com/song"
	lyrics := "Updated lyrics"
	updateData := entities.UpdateSongData{
		Band:        &band,
		Song:        &song,
		ReleaseDate: nil,
		Link:        &link,
		Lyrics:      &lyrics,
	}

	mockTM.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	mockSongRepo.On("Update", ctx, songID, updateData).Return(nil)
	mockLyricsRepo.On("Update", ctx, songID, updateData).Return(nil)

	err := useCase.Execute(ctx, songID, updateData)

	assert.NoError(t, err)
	mockTM.AssertExpectations(t)
	mockSongRepo.AssertExpectations(t)
	mockLyricsRepo.AssertExpectations(t)
}

func TestUpdateSongUseCase_Execute_SongRepoError(t *testing.T) {
	mockTM := new(MockTransactionManager)
	mockSongRepo := new(MockSongRepo)
	mockLyricsRepo := new(MockLyricsRepo)
	useCase := usecase.NewUpdateSongUseCase(mockTM, mockSongRepo, mockLyricsRepo)

	ctx := context.Background()
	songID := 123
	band := "Updated Band"
	song := "Updated Song"
	lyrics := "Updated lyrics"
	updateData := entities.UpdateSongData{
		Band:   &band,
		Song:   &song,
		Lyrics: &lyrics,
	}

	expectedError := errors.New("song repository error")

	mockTM.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).Return(expectedError)
	mockSongRepo.On("Update", ctx, songID, updateData).Return(expectedError)

	err := useCase.Execute(ctx, songID, updateData)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockTM.AssertExpectations(t)
	mockSongRepo.AssertExpectations(t)
	mockLyricsRepo.AssertNotCalled(t, "Update")
}

func TestUpdateSongUseCase_Execute_LyricsRepoError(t *testing.T) {
	mockTM := new(MockTransactionManager)
	mockSongRepo := new(MockSongRepo)
	mockLyricsRepo := new(MockLyricsRepo)
	useCase := usecase.NewUpdateSongUseCase(mockTM, mockSongRepo, mockLyricsRepo)

	ctx := context.Background()
	songID := 123
	band := "Updated Band"
	song := "Updated Song"
	lyrics := "Updated lyrics"
	updateData := entities.UpdateSongData{
		Band:   &band,
		Song:   &song,
		Lyrics: &lyrics,
	}

	expectedError := errors.New("lyrics repository error")

	mockTM.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).Return(expectedError)
	mockSongRepo.On("Update", ctx, songID, updateData).Return(nil)
	mockLyricsRepo.On("Update", ctx, songID, updateData).Return(expectedError)

	err := useCase.Execute(ctx, songID, updateData)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockTM.AssertExpectations(t)
	mockSongRepo.AssertExpectations(t)
	mockLyricsRepo.AssertExpectations(t)
}
