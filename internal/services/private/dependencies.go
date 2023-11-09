package private

import (
	"context"

	"github.com/farawaygg/wisdom/internal/wisdom"
)

type Wisdom interface {
	GetWisdoms(ctx context.Context) ([]*wisdom.Wisdom, error)
}
