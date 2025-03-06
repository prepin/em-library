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

type MockLyricsRepo struct {
	mock.Mock
}

func (m *MockLyricsRepo) Create(ctx context.Context, data entities.NewLyricsData) error {
	args := m.Called(ctx, data)
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
