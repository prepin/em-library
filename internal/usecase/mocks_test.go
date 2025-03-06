package usecase_test

import (
	"context"
	"em-library/internal/entities"

	"github.com/stretchr/testify/mock"
)

type MockTransactionManager struct {
	mock.Mock
}

func (m *MockTransactionManager) Do(ctx context.Context, f func(ctx context.Context) error) error {
	args := m.Called(ctx, f)
	if f != nil {
		f(ctx)
	}
	return args.Error(0)
}

type MockSongRepo struct {
	mock.Mock
}

func (m *MockSongRepo) Create(ctx context.Context, data entities.NewSongData) (int, error) {
	args := m.Called(ctx, data)
	return args.Int(0), args.Error(1)
}

func (m *MockSongRepo) GetList(ctx context.Context, filter entities.SongFilterData) ([]entities.SongData, error) {
	args := m.Called(ctx, filter)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]entities.SongData), args.Error(1)
}

func (m *MockSongRepo) Delete(ctx context.Context, songID int) error {
	args := m.Called(ctx, songID)
	return args.Error(0)
}

type MockLyricsRepo struct {
	mock.Mock
}

func (m *MockLyricsRepo) Create(ctx context.Context, data entities.NewLyricsData) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockLyricsRepo) Get(ctx context.Context, songID int) (entities.LyricsData, error) {
	args := m.Called(ctx, songID)

	if args.Get(0) == nil {
		return entities.LyricsData{}, args.Error(1)
	}

	return args.Get(0).(entities.LyricsData), args.Error(1)
}

func (m *MockLyricsRepo) Delete(ctx context.Context, songID int) error {
	args := m.Called(ctx, songID)
	return args.Error(0)
}

type MockSongInfoService struct {
	mock.Mock
}

func (m *MockSongInfoService) GetInfo(ctx context.Context, group, song string) (*entities.SongDetail, error) {
	args := m.Called(ctx, group, song)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.SongDetail), args.Error(1)
}
