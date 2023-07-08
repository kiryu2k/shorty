package service

import (
	"context"
	"errors"
	"time"

	"github.com/kiryu-dev/shorty/internal/libshorty/valuegen"
	"github.com/kiryu-dev/shorty/internal/model"
)

type shortenerStorage interface {
	Save(context.Context, *model.ShortURL) error
	GetAndUpdateVisits(context.Context, string) (string, error)
	GetURL(context.Context, string) (string, error)
	Delete(context.Context, string) (*model.ShortURL, error)
}

type Shortener struct {
	storage shortenerStorage
}

func NewShortener(storage shortenerStorage) *Shortener {
	return &Shortener{storage}
}

func (s *Shortener) MakeShort(ctx context.Context, url string) (string, error) {
	shortURL := &model.ShortURL{
		URL:       url,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	shortURL.Alias = valuegen.GenerateValue(shortURL.URL)
	urlFromStorage, err := s.storage.GetURL(ctx, shortURL.Alias)
	if errors.Is(err, model.ErrURLNotFound) {
		if err := s.storage.Save(ctx, shortURL); err != nil {
			return "", err
		}
		return shortURL.Alias, nil
	}
	if err != nil {
		return "", err
	}
	if url != urlFromStorage {
		return "", model.ErrCollision
	}
	return shortURL.Alias, nil
}

func (s *Shortener) GetURL(ctx context.Context, alias string) (string, error) {
	return s.storage.GetAndUpdateVisits(ctx, alias)
}
