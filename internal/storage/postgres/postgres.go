package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/serenize/snaker"

	"github.com/farawaygg/wisdom/internal/storage"
)

type Storage struct {
	db *sqlx.DB

	stInsertWisdom *sqlx.NamedStmt
	stGetWisdoms   *sqlx.Stmt
}

func New(db *sqlx.DB) (*Storage, error) {
	db.MapperFunc(snaker.CamelToSnake)
	var (
		s   = Storage{db: db}
		err error
	)

	s.stInsertWisdom, err = db.PrepareNamed(`
		INSERT INTO wisdoms(id, value, author_id, created_at, updated_at)
		VALUES(:id, :value, :author_id, :created_at, :updated_at)
	`)
	if err != nil {
		return nil, errors.WithMessage(err, "db.PrepareNamed(stInsertWisdom)")
	}

	s.stGetWisdoms, err = db.Preparex(`
		SELECT id, value, author_id, created_at, updated_at
		FROM wisdoms
	`)
	if err != nil {
		return nil, errors.WithMessage(err, "db.Preparex(stGetWisdoms)")
	}
	return &s, nil
}

func (s *Storage) GetWisdoms(ctx context.Context, iter storage.WisdomIterFunc) error {
	if err := s.queryWisdoms(ctx, iter, s.stGetWisdoms); err != nil {
		return errors.WithMessage(err, "queryWisdoms")
	}
	return nil
}

func (s *Storage) CreateWisdom(ctx context.Context, wisdom storage.Wisdom) error {
	result, err := s.stInsertWisdom.ExecContext(ctx, wisdom)
	if err != nil {
		return errors.WithMessage(storageError{err}, "stInsertWisdom.ExecContext")
	}

	return shouldAffectRows(result)
}

func (s *Storage) queryWisdoms(
	ctx context.Context,
	iter storage.WisdomIterFunc,
	stmt *sqlx.Stmt, args ...any) error {

	rows, err := stmt.QueryxContext(ctx, args...)
	if err != nil {
		return storageError{err}
	}

	defer rows.Close()

	for rows.Next() {
		var row storage.Wisdom
		if err := rows.StructScan(&row); err != nil {
			return errors.WithMessage(err, "scan wallet account record")
		}

		if err := iter(row); err != nil {
			if errors.Is(err, storage.ErrStopIteration) {
				break
			}

			return errors.WithMessage(err, "iter")
		}
	}

	return rows.Err()
}

type storageError struct {
	error
}

func (s storageError) Unwrap() error {
	causeErr := errors.Cause(s.error)
	if errors.Is(causeErr, sql.ErrNoRows) {
		return storage.ErrNotFound
	}

	return causeErr
}

func shouldAffectRows(rs sql.Result) error {
	n, err := rs.RowsAffected()
	if err != nil {
		return errors.WithMessage(err, "RowsAffected")
	}

	if n <= 0 {
		return storage.ErrNotFound
	}

	return nil
}
