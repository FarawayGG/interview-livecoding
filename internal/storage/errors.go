package storage

import (
	"github.com/pkg/errors"
)

var (
	ErrNotFound = errors.New("not found")

	ErrStopIteration = errors.New("stop iteration")
)
