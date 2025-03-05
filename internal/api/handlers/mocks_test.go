package handlers_test

import (
	"context"
	"em-library/internal/entities"

	"github.com/stretchr/testify/mock"
)

// MockLogger implements the Logger interface for testing
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockLogger) Info(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockLogger) Error(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockLogger) Warn(msg string, args ...interface{}) {
	m.Called(msg, args)
}

// MockCreateSongUseCase mocks the CreateSongUseCase for testing
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

// MockUseCases combines all mock use cases
type MockUseCases struct {
	CreateSong *MockCreateSongUseCase
}
