package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kiryu-dev/shorty/internal/config"
	"github.com/kiryu-dev/shorty/internal/model"
	"github.com/kiryu-dev/shorty/internal/storage"
	_ "github.com/lib/pq"
)

type Storage struct {
	*storage.Queries
	db *sql.DB
}

func New(cfg *config.DB) (*Storage, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, fmt.Errorf("invalid connection to postgres: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cannot get access to postgtes: %w", err)
	}
	return &Storage{
		db:      db,
		Queries: storage.New(db),
	}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) execTx(ctx context.Context, fn func(*storage.Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := storage.New(tx)
	if err := fn(q); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("transaction error: %w, rollback error: %w", err, rollbackErr)
		}
		return err
	}
	return tx.Commit()
}

func (s *Storage) Save(ctx context.Context, shorty *model.ShortURL) error {
	const (
		op    = "postgres.Storage.Save"
		query = `
INSERT INTO urls (url, alias, created_at, updated_at)
VALUES ($1, $2, $3, $4) RETURNING id;
		`
	)
	err := s.execTx(ctx, func(q *storage.Queries) error {
		err := q.DB().QueryRowContext(ctx, query, shorty.URL, shorty.Alias,
			shorty.CreatedAt, shorty.UpdatedAt).Scan(&shorty.ID)
		return err
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) GetAndUpdateVisits(ctx context.Context, alias string) (string, error) {
	const (
		op    = "postgres.Storage.GetAndUpdateVisits"
		query = `SELECT url FROM urls WHERE alias = $1;`
	)
	var url string
	err := s.execTx(ctx, func(_ *storage.Queries) error {
		var err error
		url, err = s.GetURL(ctx, alias)
		if err != nil {
			return err
		}
		return s.incrementVisits(ctx, alias)
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return url, nil
}

func (s *Storage) incrementVisits(ctx context.Context, alias string) error {
	const (
		op    = "postgres.Storage.incrementVisits"
		query = `
UPDATE urls
SET visits = visits + 1,
	updated_at = NOW()
WHERE alias = $1;
		`
	)
	res, err := s.db.ExecContext(ctx, query, alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if count == 0 {
		return fmt.Errorf("%s: %w", op, model.ErrURLNotFound)
	}
	return nil
}

func (s *Storage) GetURL(ctx context.Context, alias string) (string, error) {
	const (
		op    = "postgres.Storage.GetURL"
		query = `SELECT url FROM urls WHERE alias = $1;`
	)
	var url string
	err := s.db.QueryRowContext(ctx, query, alias).Scan(&url)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("%s: %w", op, model.ErrURLNotFound)
	}
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return url, nil
}

func (s *Storage) Delete(ctx context.Context, alias string) (*model.ShortURL, error) {
	const (
		op          = "postgres.Storage.Delete"
		findQuery   = `SELECT id, url, visits, created_at FROM urls WHERE alias = $1;`
		deleteQuery = `DELETE FROM urls WHERE alias = $1;`
	)
	del := new(model.ShortURL)
	err := s.execTx(ctx, func(q *storage.Queries) error {
		err := q.DB().QueryRowContext(ctx, findQuery, alias).
			Scan(&del.ID, &del.URL, &del.Visits, &del.CreatedAt)
		if err == sql.ErrNoRows {
			return model.ErrURLNotFound
		}
		if err != nil {
			return err
		}
		_, err = q.DB().ExecContext(ctx, deleteQuery, alias)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return del, nil
}
