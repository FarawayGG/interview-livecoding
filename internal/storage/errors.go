package storage

import "errors"

var (
	ErrNotFound = errors.New("not found")

	ErrStopIteration = errors.New("stop iteration")
)
