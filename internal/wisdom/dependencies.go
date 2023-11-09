package wisdom

import (
	"context"

	"github.com/farawaygg/wisdom/internal/storage"
)

type Storage interface {
	GetWisdoms(ctx context.Context, iter storage.WisdomIterFunc) error
}
