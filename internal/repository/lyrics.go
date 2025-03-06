package repository

import (
	"context"
	"em-library/config"
	"em-library/internal/entities"
	"em-library/internal/errs"
	"em-library/pkg/database"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dm"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
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

func (r *PGLyricsRepository) Get(ctx context.Context, songID int) (entities.LyricsData, error) {
	stmt := psql.Select(
		sm.Columns("content"),
		sm.From("lyrics"),
		sm.Where(psql.Quote("song_id").EQ(psql.Arg(songID))),
	)

	query, args := stmt.MustBuild(ctx)
	r.logger.Debug("executing select lyrics query", "query", query, "args", args)

	var content string
	err := r.db.Conn(ctx).QueryRow(ctx, query, args...).Scan(&content)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.LyricsData{}, fmt.Errorf("%w song lyrics not found", errs.ErrNotFound)
		}
		return entities.LyricsData{}, err
	}
	r.logger.Debug("lyrics queried successfully", "song_id", songID)

	return entities.LyricsData{
		SongID:  songID,
		Content: content,
	}, nil
}

func (r *PGLyricsRepository) Delete(ctx context.Context, songID int) error {

	stmt := psql.Delete(
		dm.From("lyrics"),
		dm.Where(psql.Quote("song_id").EQ(psql.Arg(songID))),
	)

	query, args := stmt.MustBuild(ctx)
	r.logger.Debug("executing delete lyrics query", "query", query, "args", args)

	_, err := r.db.Conn(ctx).Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	r.logger.Debug("lyrics deleted successfully", "song_id", songID)

	return nil
}
