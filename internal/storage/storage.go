package storage

import "fmt"

var (
	ErrURLNotFound = fmt.Errorf("url is not found")
	ErrURLExists   = fmt.Errorf("url exists")
)
