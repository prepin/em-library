package repository

import (
	"context"
	"em-library/internal/entities"
)

type PGSongRepository struct {
	// TODO: add logger
}

func NewPGSongRepository() *PGSongRepository {
	return &PGSongRepository{}
}

func (r *PGSongRepository) Create(ctx context.Context, data entities.NewSongData) (int, error) {
	return 124, nil
}
