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

func TestGetSongListUseCase_Execute_Success(t *testing.T) {
	mockSongRepo := new(MockSongRepo)
	useCase := usecase.NewGetSongListUseCase(mockSongRepo)

	ctx := context.Background()
	releaseDate := time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC)

	mockSongs := []entities.SongData{
		{ID: 1, Band: "Band1", Song: "Song1", ReleaseDate: releaseDate, Link: "link1"},
		{ID: 2, Band: "Band2", Song: "Song2", ReleaseDate: releaseDate, Link: "link2"},
	}

	band := "Band"
	limit := 10
	offset := 0
	filter := entities.SongFilterData{
		Band:   &band,
		Limit:  &limit,
		Offset: &offset,
	}

	expectedFilter := entities.SongFilterData{
		Band:   &band,
		Limit:  &limit,
		Offset: &offset,
	}

	mockSongRepo.On("GetList", ctx, expectedFilter).Return(mockSongs, nil)

	result, err := useCase.Execute(ctx, filter)

	assert.NoError(t, err)
	assert.Equal(t, mockSongs, result)
	assert.Len(t, result, 2)
	mockSongRepo.AssertExpectations(t)
}

func TestGetSongListUseCase_Execute_WithDefaultLimitAndOffset(t *testing.T) {
	mockSongRepo := new(MockSongRepo)
	useCase := usecase.NewGetSongListUseCase(mockSongRepo)

	ctx := context.Background()
	releaseDate := time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC)

	mockSongs := []entities.SongData{
		{ID: 1, Band: "Band1", Song: "Song1", ReleaseDate: releaseDate, Link: "link1"},
		{ID: 2, Band: "Band2", Song: "Song2", ReleaseDate: releaseDate, Link: "link2"},
	}

	filter := entities.SongFilterData{}

	expectedFilter := entities.SongFilterData{}
	expectedLimit := 50
	expectedOffset := 0
	expectedFilter.Limit = &expectedLimit
	expectedFilter.Offset = &expectedOffset

	mockSongRepo.On("GetList", ctx, mock.MatchedBy(func(f entities.SongFilterData) bool {
		return *f.Limit == 50 && *f.Offset == 0
	})).Return(mockSongs, nil)

	result, err := useCase.Execute(ctx, filter)

	assert.NoError(t, err)
	assert.Equal(t, mockSongs, result)
	assert.Len(t, result, 2)
	mockSongRepo.AssertExpectations(t)
}

func TestGetSongListUseCase_Execute_WithDefaultLimit(t *testing.T) {
	mockSongRepo := new(MockSongRepo)
	useCase := usecase.NewGetSongListUseCase(mockSongRepo)

	ctx := context.Background()
	releaseDate := time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC)

	mockSongs := []entities.SongData{
		{ID: 1, Band: "Band1", Song: "Song1", ReleaseDate: releaseDate, Link: "link1"},
		{ID: 2, Band: "Band2", Song: "Song2", ReleaseDate: releaseDate, Link: "link2"},
	}

	offset := 5
	filter := entities.SongFilterData{
		Offset: &offset,
	}

	mockSongRepo.On("GetList", ctx, mock.MatchedBy(func(f entities.SongFilterData) bool {
		return *f.Limit == 50 && *f.Offset == 5
	})).Return(mockSongs, nil)

	result, err := useCase.Execute(ctx, filter)

	assert.NoError(t, err)
	assert.Equal(t, mockSongs, result)
	assert.Len(t, result, 2)
	mockSongRepo.AssertExpectations(t)
}

func TestGetSongListUseCase_Execute_WithDefaultOffset(t *testing.T) {
	mockSongRepo := new(MockSongRepo)
	useCase := usecase.NewGetSongListUseCase(mockSongRepo)

	ctx := context.Background()
	releaseDate := time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC)

	mockSongs := []entities.SongData{
		{ID: 1, Band: "Band1", Song: "Song1", ReleaseDate: releaseDate, Link: "link1"},
		{ID: 2, Band: "Band2", Song: "Song2", ReleaseDate: releaseDate, Link: "link2"},
	}

	limit := 25
	filter := entities.SongFilterData{
		Limit: &limit,
	}

	mockSongRepo.On("GetList", ctx, mock.MatchedBy(func(f entities.SongFilterData) bool {
		return *f.Limit == 25 && *f.Offset == 0
	})).Return(mockSongs, nil)

	result, err := useCase.Execute(ctx, filter)

	assert.NoError(t, err)
	assert.Equal(t, mockSongs, result)
	assert.Len(t, result, 2)
	mockSongRepo.AssertExpectations(t)
}

func TestGetSongListUseCase_Execute_WithComplexFilter(t *testing.T) {
	mockSongRepo := new(MockSongRepo)
	useCase := usecase.NewGetSongListUseCase(mockSongRepo)

	ctx := context.Background()
	releaseDate := time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC)
	releaseDateFrom := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	releaseDateTo := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)

	mockSongs := []entities.SongData{
		{ID: 1, Band: "Band1", Song: "Song1", ReleaseDate: releaseDate, Link: "link1"},
	}

	band := "Band1"
	song := "Song1"
	id := 1
	limit := 10
	offset := 0
	filter := entities.SongFilterData{
		ID:              &id,
		Band:            &band,
		Song:            &song,
		ReleaseDateFrom: &releaseDateFrom,
		ReleaseDateTo:   &releaseDateTo,
		Limit:           &limit,
		Offset:          &offset,
	}

	mockSongRepo.On("GetList", ctx, filter).Return(mockSongs, nil)

	result, err := useCase.Execute(ctx, filter)

	assert.NoError(t, err)
	assert.Equal(t, mockSongs, result)
	assert.Len(t, result, 1)
	mockSongRepo.AssertExpectations(t)
}

func TestGetSongListUseCase_Execute_EmptyResult(t *testing.T) {
	mockSongRepo := new(MockSongRepo)
	useCase := usecase.NewGetSongListUseCase(mockSongRepo)

	ctx := context.Background()
	mockSongs := []entities.SongData{}

	band := "NonexistentBand"
	limit := 10
	offset := 0
	filter := entities.SongFilterData{
		Band:   &band,
		Limit:  &limit,
		Offset: &offset,
	}

	mockSongRepo.On("GetList", ctx, filter).Return(mockSongs, nil)

	result, err := useCase.Execute(ctx, filter)

	assert.NoError(t, err)
	assert.Empty(t, result)
	mockSongRepo.AssertExpectations(t)
}

func TestGetSongListUseCase_Execute_RepoError(t *testing.T) {
	mockSongRepo := new(MockSongRepo)
	useCase := usecase.NewGetSongListUseCase(mockSongRepo)

	ctx := context.Background()
	expectedError := errors.New("database error")

	limit := 10
	offset := 0
	filter := entities.SongFilterData{
		Limit:  &limit,
		Offset: &offset,
	}

	mockSongRepo.On("GetList", ctx, filter).Return(nil, expectedError)

	result, err := useCase.Execute(ctx, filter)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)
	mockSongRepo.AssertExpectations(t)
}
