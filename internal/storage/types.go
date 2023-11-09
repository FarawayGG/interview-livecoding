package storage

import (
	"time"

	"github.com/google/uuid"
)

type Wisdom struct {
	ID        uuid.UUID
	Value     string
	AuthorID  uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}
