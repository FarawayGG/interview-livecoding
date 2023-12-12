package wisdom

import (
	"context"
	"fmt"

	"github.com/farawaygg/wisdom/internal/storage"
)

type Repo struct {
	storage Storage
}

func New(s Storage) *Repo {
	return &Repo{
		storage: s,
	}
}

func (r *Repo) GetWisdoms(ctx context.Context) ([]*Wisdom, error) {
	var wisdoms []*Wisdom
	if err := r.storage.GetWisdoms(ctx, func(w storage.Wisdom) error {
		wisdoms = append(wisdoms, &Wisdom{
			Value:     w.Value,
			CreatedAt: w.CreatedAt,
		})
		return nil
	}); err != nil {
		return nil, fmt.Errorf("storage.GetWisdoms: %w", err)
	}

	return wisdoms, nil
}
