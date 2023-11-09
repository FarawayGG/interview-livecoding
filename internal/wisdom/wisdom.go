package wisdom

import (
	"context"

	"github.com/pkg/errors"

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
		return nil, errors.WithMessage(err, "storage.GetWisdoms")
	}

	return wisdoms, nil
}
