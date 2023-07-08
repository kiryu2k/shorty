package model

import (
	"fmt"
	"time"
)

type ShortURL struct {
	ID        int64     `json:"id"`
	URL       string    `json:"url"`
	Alias     string    `json:"alias"`
	Visits    uint64    `json:"visits"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
}

var (
	ErrURLNotFound = fmt.Errorf("url is not found")
	ErrURLExists   = fmt.Errorf("url exists")
	ErrCollision   = fmt.Errorf("unexpected collision when generating short url")
)
