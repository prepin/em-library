package usecase_test

import (
	"context"
	"em-library/internal/usecase"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteSongUseCase_Execute_Success(t *testing.T) {
	mockTransactionManager := new(MockTransactionManager)
	mockSongRepo := new(MockSongRepo)
	mockLyricsRepo := new(MockLyricsRepo)
	useCase := usecase.NewDeleteSongUseCase(mockTransactionManager, mockSongRepo, mockLyricsRepo)

	ctx := context.Background()
	songID := 1

	mockTransactionManager.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	mockLyricsRepo.On("Delete", ctx, songID).Return(nil)
	mockSongRepo.On("Delete", ctx, songID).Return(nil)

	err := useCase.Execute(ctx, songID)

	assert.NoError(t, err)
	mockTransactionManager.AssertExpectations(t)
	mockLyricsRepo.AssertExpectations(t)
	mockSongRepo.AssertExpectations(t)
}

func TestDeleteSongUseCase_Execute_LyricsDeleteError(t *testing.T) {
	mockTransactionManager := new(MockTransactionManager)
	mockSongRepo := new(MockSongRepo)
	mockLyricsRepo := new(MockLyricsRepo)
	useCase := usecase.NewDeleteSongUseCase(mockTransactionManager, mockSongRepo, mockLyricsRepo)

	ctx := context.Background()
	songID := 1
	expectedError := errors.New("lyrics delete error")

	mockLyricsRepo.On("Delete", ctx, songID).Return(expectedError)
	mockTransactionManager.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(context.Context) error)
		fn(ctx)
	}).Return(expectedError)

	err := useCase.Execute(ctx, songID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockTransactionManager.AssertExpectations(t)
	mockLyricsRepo.AssertExpectations(t)
	mockSongRepo.AssertNotCalled(t, "Delete", ctx, songID)
}

func TestDeleteSongUseCase_Execute_SongDeleteError(t *testing.T) {
	mockTransactionManager := new(MockTransactionManager)
	mockSongRepo := new(MockSongRepo)
	mockLyricsRepo := new(MockLyricsRepo)
	useCase := usecase.NewDeleteSongUseCase(mockTransactionManager, mockSongRepo, mockLyricsRepo)

	ctx := context.Background()
	songID := 1
	expectedError := errors.New("song delete error")

	mockLyricsRepo.On("Delete", ctx, songID).Return(nil)
	mockSongRepo.On("Delete", ctx, songID).Return(expectedError)
	mockTransactionManager.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(context.Context) error)
		fn(ctx)
	}).Return(expectedError)

	err := useCase.Execute(ctx, songID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockTransactionManager.AssertExpectations(t)
	mockLyricsRepo.AssertExpectations(t)
	mockSongRepo.AssertExpectations(t)
}
