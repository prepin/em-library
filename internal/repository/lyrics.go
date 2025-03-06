package repository

import (
	"context"
	"em-library/config"
	"em-library/internal/entities"
	"em-library/pkg/database"

	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/im"
)

type PGLyricsRepository struct {
	db     *database.Database
	logger config.Logger
}

func NewPGLyricsRepository(db *database.Database, l config.Logger) *PGLyricsRepository {
	return &PGLyricsRepository{
		db:     db,
		logger: l,
	}
}

func (r *PGLyricsRepository) Create(ctx context.Context, data entities.NewLyricsData) error {
	stmt := psql.Insert(
		im.Into("lyrics", "song_id", "content"),
		im.Values(
			psql.Arg(data.SongID),
			psql.Arg(data.Content),
		),
	)

	query, args := stmt.MustBuild(ctx)
	r.logger.Debug("executing insert lyrics query", "query", query, "args", args)

	_, err := r.db.Conn(ctx).Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	r.logger.Debug("lyrics inserted successfully", "song_id", data.SongID)

	return nil
}
