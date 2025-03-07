package repository

import (
	"context"
	"em-library/config"
	"em-library/internal/entities"
	"em-library/internal/errs"
	"em-library/pkg/database"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dm"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/dialect/psql/um"
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

func (r *PGSongRepository) GetList(
	ctx context.Context,
	filter entities.SongFilterData,
) ([]entities.SongData, error) {

	stmt := psql.Select(
		sm.Columns("id", "band", "song", "release_date", "link"),
		sm.From("songs"),
		sm.OrderBy("release_date"),
	)

	if filter.ID != nil {
		stmt.Apply(sm.Where(psql.Quote("id").EQ(psql.Arg(*filter.ID))))
	}

	if filter.Band != nil {
		stmt.Apply(sm.Where(psql.Quote("band").EQ(psql.Arg(*filter.Band))))
	}

	if filter.Song != nil {
		stmt.Apply(sm.Where(psql.Quote("song").EQ(psql.Arg(*filter.Song))))
	}

	if filter.ReleaseDateFrom != nil {
		stmt.Apply(sm.Where(psql.Quote("release_date").GTE(psql.Arg(*filter.ReleaseDateFrom))))
	}

	if filter.ReleaseDateTo != nil {
		stmt.Apply(sm.Where(psql.Quote("release_date").LTE(psql.Arg(*filter.ReleaseDateTo))))
	}

	if filter.Offset != nil {
		stmt.Apply(sm.Offset(*filter.Offset))
	}

	if filter.Limit != nil {
		stmt.Apply(sm.Limit(*filter.Limit))
	}

	query, args := stmt.MustBuild(ctx)
	r.logger.Debug("executing select song list query", "query", query, "args", args)

	rows, err := r.db.Conn(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	songs, err := pgx.CollectRows(rows, pgx.RowToStructByName[entities.SongData])
	if err != nil {
		return nil, err
	}

	if len(songs) == 0 {
		return nil, fmt.Errorf("%w songs not found", errs.ErrNotFound)
	}

	r.logger.Debug("Successfully queried songs", "count", len(songs))
	return songs, nil
}

func (r *PGSongRepository) Update(ctx context.Context, songID int, data entities.UpdateSongData) error {

	nothingToUpdate := true

	stmt := psql.Update(
		um.Table("songs"),
		um.SetCol("updated_at").ToArg(time.Now()),
		um.Where(psql.Quote("id").EQ(psql.Arg(songID))),
	)

	if data.Band != nil {
		stmt.Apply(
			um.SetCol("band").ToArg(*data.Band),
		)
		nothingToUpdate = false
	}

	if data.Song != nil {
		stmt.Apply(
			um.SetCol("song").ToArg(*data.Song),
		)
		nothingToUpdate = false
	}

	if data.ReleaseDate != nil {
		stmt.Apply(
			um.SetCol("release_date").ToArg(*data.ReleaseDate),
		)
		nothingToUpdate = false
	}

	if data.Link != nil {
		stmt.Apply(
			um.SetCol("link").ToArg(*data.Link),
		)
		nothingToUpdate = false
	}

	if nothingToUpdate {
		return nil
	}

	query, args := stmt.MustBuild(ctx)
	r.logger.Debug("executing update song query", "query", query, "args", args)

	ct, err := r.db.Conn(ctx).Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("%w no song rows updated", errs.ErrNotFound)
	}

	r.logger.Debug("song updated successfully", "id", songID)

	return nil
}

func (r *PGSongRepository) Delete(ctx context.Context, songID int) error {
	stmt := psql.Delete(
		dm.From("songs"),
		dm.Where(psql.Quote("id").EQ(psql.Arg(songID))),
	)

	query, args := stmt.MustBuild(ctx)
	r.logger.Debug("executing delete song query", "query", query, "args", args)

	_, err := r.db.Conn(ctx).Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	r.logger.Debug("song deleted successfully", "id", songID)

	return nil
}
