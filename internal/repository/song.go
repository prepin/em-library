package repository

import (
	"context"
	"em-library/config"
	"em-library/internal/entities"
	"em-library/internal/errs"
	"em-library/pkg/database"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/im"
)

type PGSongRepository struct {
	db     *database.Database
	logger config.Logger
}

func NewPGSongRepository(db *database.Database, l config.Logger) *PGSongRepository {
	return &PGSongRepository{
		db:     db,
		logger: l,
	}
}

const (
	PG_ERROR_EXISTS       = "23505"
	SONG_BAND_UNIQ_CONSTR = "unique_band_song"
)

func (r *PGSongRepository) Create(ctx context.Context, data entities.NewSongData) (int, error) {
	stmt := psql.Insert(
		im.Into("songs", "band", "song", "release_date", "link"),
		im.Values(
			psql.Arg(data.Band),
			psql.Arg(data.Song),
			psql.Arg(data.ReleaseDate),
			psql.Arg(data.Link),
		),
		im.Returning("id"),
	)

	query, args := stmt.MustBuild(ctx)

	var id int
	r.logger.Debug("executing insert song query", "query", query, "args", args)
	err := r.db.Conn(ctx).QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok && pgErr.Code == PG_ERROR_EXISTS && pgErr.ConstraintName == SONG_BAND_UNIQ_CONSTR {
			return 0, fmt.Errorf("%w song '%s' by '%s' already exists", errs.ErrAlreadyExists, data.Song, data.Band)
		}
		return 0, err
	}
	r.logger.Debug("song inserted successfully", "id", id)
	return id, nil
}
