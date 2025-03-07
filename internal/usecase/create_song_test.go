package usecase_test

import (
	"context"
	"em-library/internal/entities"
	"em-library/internal/usecase"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateSongUseCase_Execute_Success(t *testing.T) {
	mockTM := new(MockTransactionManager)
	mockSongRepo := new(MockSongRepo)
	mockLyricsRepo := new(MockLyricsRepo)
	mockInfoService := new(MockSongInfoService)
	useCase := usecase.NewCreateSongUseCase(mockTM, mockSongRepo, mockLyricsRepo, mockInfoService)

	ctx := context.Background()
	expectedID := 123
	releaseDate := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)

	inputData := entities.NewSongData{
		Band: "Test Group",
		Song: "Test Song",
	}

	songDetail := &entities.SongDetail{
		ReleaseDate: releaseDate,
		Link:        "https://example.com/song",
		Lyrics:        "Song lyrics",
	}

	expectedData := inputData
	expectedData.ReleaseDate = releaseDate
	expectedData.Link = "https://example.com/song"

	mockInfoService.On("GetInfo", ctx, inputData.Band, inputData.Song).Return(songDetail, nil)
	mockTM.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	mockSongRepo.On("Create", ctx, expectedData).Return(expectedID, nil)
	mockLyricsRepo.On("Create", ctx, entities.NewLyricsData{
		SongID:  expectedID,
		Content: songDetail.Lyrics,
	}).Return(nil)

	result, err := useCase.Execute(ctx, inputData)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedID, result.ID)
	assert.Equal(t, inputData.Band, result.Band)
	assert.Equal(t, inputData.Song, result.Song)
	assert.Equal(t, songDetail.ReleaseDate, result.ReleaseDate)
	assert.Equal(t, songDetail.Link, result.Link)

	mockInfoService.AssertExpectations(t)
	mockTM.AssertExpectations(t)
	mockSongRepo.AssertExpectations(t)
	mockLyricsRepo.AssertExpectations(t)
}

func TestCreateSongUseCase_Execute_InfoServiceError(t *testing.T) {
	mockTM := new(MockTransactionManager)
	mockSongRepo := new(MockSongRepo)
	mockLyricsRepo := new(MockLyricsRepo)
	mockInfoService := new(MockSongInfoService)
	useCase := usecase.NewCreateSongUseCase(mockTM, mockSongRepo, mockLyricsRepo, mockInfoService)

	ctx := context.Background()
	inputData := entities.NewSongData{
		Band: "Test Group",
		Song: "Test Song",
	}

	expectedError := errors.New("info service error")
	mockInfoService.On("GetInfo", ctx, inputData.Band, inputData.Song).Return(nil, expectedError)

	result, err := useCase.Execute(ctx, inputData)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)

	mockInfoService.AssertExpectations(t)
	mockTM.AssertNotCalled(t, "Do")
	mockSongRepo.AssertNotCalled(t, "Create")
	mockLyricsRepo.AssertNotCalled(t, "Create")
}

func TestCreateSongUseCase_Execute_SongRepoError(t *testing.T) {
	mockTM := new(MockTransactionManager)
	mockSongRepo := new(MockSongRepo)
	mockLyricsRepo := new(MockLyricsRepo)
	mockInfoService := new(MockSongInfoService)
	useCase := usecase.NewCreateSongUseCase(mockTM, mockSongRepo, mockLyricsRepo, mockInfoService)

	ctx := context.Background()
	inputData := entities.NewSongData{
		Band: "Test Group",
		Song: "Test Song",
	}

	releaseDate := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	songDetail := &entities.SongDetail{
		ReleaseDate: releaseDate,
		Link:        "https://example.com/song",
		Lyrics:        "Song lyrics",
	}

	expectedData := inputData
	expectedData.ReleaseDate = releaseDate
	expectedData.Link = "https://example.com/song"

	expectedError := errors.New("repository error")

	mockInfoService.On("GetInfo", ctx, inputData.Band, inputData.Song).Return(songDetail, nil)
	mockTM.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).Return(expectedError)
	mockSongRepo.On("Create", ctx, expectedData).Return(0, expectedError)

	result, err := useCase.Execute(ctx, inputData)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)

	mockInfoService.AssertExpectations(t)
	mockTM.AssertExpectations(t)
}

func TestCreateSongUseCase_Execute_LyricsRepoError(t *testing.T) {
	mockTM := new(MockTransactionManager)
	mockSongRepo := new(MockSongRepo)
	mockLyricsRepo := new(MockLyricsRepo)
	mockInfoService := new(MockSongInfoService)
	useCase := usecase.NewCreateSongUseCase(mockTM, mockSongRepo, mockLyricsRepo, mockInfoService)

	ctx := context.Background()
	expectedID := 123
	releaseDate := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)

	inputData := entities.NewSongData{
		Band: "Test Group",
		Song: "Test Song",
	}

	songDetail := &entities.SongDetail{
		ReleaseDate: releaseDate,
		Link:        "https://example.com/song",
		Lyrics:        "Song lyrics",
	}

	expectedData := inputData
	expectedData.ReleaseDate = releaseDate
	expectedData.Link = "https://example.com/song"

	expectedError := errors.New("lyrics repository error")

	mockInfoService.On("GetInfo", ctx, inputData.Band, inputData.Song).Return(songDetail, nil)
	mockTM.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).Return(expectedError)
	mockSongRepo.On("Create", ctx, expectedData).Return(expectedID, nil)
	mockLyricsRepo.On("Create", ctx, entities.NewLyricsData{
		SongID:  expectedID,
		Content: songDetail.Lyrics,
	}).Return(expectedError)

	result, err := useCase.Execute(ctx, inputData)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)

	mockInfoService.AssertExpectations(t)
	mockTM.AssertExpectations(t)
	mockSongRepo.AssertExpectations(t)
	mockLyricsRepo.AssertExpectations(t)
}
