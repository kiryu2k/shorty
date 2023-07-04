package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kiryu-dev/shorty/internal/config"
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
		return nil, fmt.Errorf("invalid connection to postgres: %s", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cannot get access to postgtes: %s", err)
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
			return fmt.Errorf("transaction error: %s, rollback error: %s", err, rollbackErr)
		}
		return err
	}
	return tx.Commit()
}
