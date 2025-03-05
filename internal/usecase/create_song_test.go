package usecase_test

import (
	"context"
	"em-library/internal/entities"
	"em-library/internal/usecase"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Happy path, всё отработало как положено
func TestCreateSongUseCase_Execute_Success(t *testing.T) {
	mockRepo := new(MockSongRepo)
	mockInfoService := new(MockSongInfoService)
	useCase := usecase.NewCreateSongUseCase(mockRepo, mockInfoService)

	ctx := context.Background()
	inputData := entities.NewSongData{
		Group: "Test Group",
		Song:  "Test Song",
	}

	expectedID := 123
	releaseDate := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	songDetail := &entities.SongDetail{
		ReleaseDate: releaseDate,
		Link:        "https://example.com/song",
	}

	mockInfoService.On("GetInfo", ctx, inputData.Group, inputData.Song).Return(songDetail, nil)
	mockRepo.On("Create", ctx, inputData).Return(expectedID, nil)

	result, err := useCase.Execute(ctx, inputData)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedID, result.ID)
	assert.Equal(t, inputData.Group, result.Group)
	assert.Equal(t, inputData.Song, result.Song)
	assert.Equal(t, songDetail.ReleaseDate, result.ReleaseDate)
	assert.Equal(t, songDetail.Link, result.Link)

	mockInfoService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// Сервис с данными песен вернул ошибку
func TestCreateSongUseCase_Execute_InfoServiceError(t *testing.T) {
	mockRepo := new(MockSongRepo)
	mockInfoService := new(MockSongInfoService)
	useCase := usecase.NewCreateSongUseCase(mockRepo, mockInfoService)

	ctx := context.Background()
	inputData := entities.NewSongData{
		Group: "Test Group",
		Song:  "Test Song",
	}

	expectedError := errors.New("info service error")
	mockInfoService.On("GetInfo", ctx, inputData.Group, inputData.Song).Return(nil, expectedError)

	result, err := useCase.Execute(ctx, inputData)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)

	mockInfoService.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Create")
}

// Репозиторий вернул ошибку при вставке
func TestCreateSongUseCase_Execute_RepoError(t *testing.T) {
	mockRepo := new(MockSongRepo)
	mockInfoService := new(MockSongInfoService)
	useCase := usecase.NewCreateSongUseCase(mockRepo, mockInfoService)

	ctx := context.Background()
	inputData := entities.NewSongData{
		Group: "Test Group",
		Song:  "Test Song",
	}

	releaseDate := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	songDetail := &entities.SongDetail{
		ReleaseDate: releaseDate,
		Link:        "https://example.com/song",
	}

	expectedError := errors.New("repository error")

	mockInfoService.On("GetInfo", ctx, inputData.Group, inputData.Song).Return(songDetail, nil)
	mockRepo.On("Create", ctx, inputData).Return(0, expectedError)

	result, err := useCase.Execute(ctx, inputData)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)

	mockInfoService.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}
