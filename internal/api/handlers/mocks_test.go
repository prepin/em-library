package handlers_test

import (
	"context"
	"em-library/internal/entities"

	"github.com/stretchr/testify/mock"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(msg string, args ...any) {
	m.Called(msg, args)
}

func (m *MockLogger) Info(msg string, args ...any) {
	m.Called(msg, args)
}

func (m *MockLogger) Error(msg string, args ...any) {
	m.Called(msg, args)
}

func (m *MockLogger) Warn(msg string, args ...any) {
	m.Called(msg, args)
}

type MockCreateSongUseCase struct {
	mock.Mock
}

func (m *MockCreateSongUseCase) Execute(ctx context.Context, data entities.NewSongData) (*entities.SongData, error) {
	args := m.Called(ctx, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.SongData), args.Error(1)
}

type MockUseCases struct {
	CreateSong *MockCreateSongUseCase
}

type MockGetSongListUseCase struct {
	mock.Mock
}

func (m *MockGetSongListUseCase) Execute(ctx context.Context, filter entities.SongFilterData) ([]entities.SongData, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.SongData), args.Error(1)
}

type MockGetSongLyricsUseCase struct {
	mock.Mock
}

func (m *MockGetSongLyricsUseCase) Execute(ctx context.Context, songID int, filter entities.LyricsFilterData) ([]entities.LyricsVerseData, error) {
	args := m.Called(ctx, songID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.LyricsVerseData), args.Error(1)
}
