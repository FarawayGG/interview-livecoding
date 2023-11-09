package private

import (
	"context"

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/farawaygg/wisdom/internal/wisdom"
	api "github.com/farawaygg/wisdom/pkg/wisdom"
)

func (s Service) GetWisdoms(ctx context.Context, _ *api.GetWisdoms_Request) (*api.GetWisdoms_Response, error) {
	wisdoms, err := s.wisdoms.GetWisdoms(ctx)
	if err != nil {
		return nil, err
	}

	return &api.GetWisdoms_Response{
		Wisdoms: lo.Map(wisdoms, func(w *wisdom.Wisdom, _ int) *api.Wisdom {
			return &api.Wisdom{
				Value:     w.Value,
				CreatedAt: timestamppb.New(w.CreatedAt),
			}
		}),
	}, nil
}
