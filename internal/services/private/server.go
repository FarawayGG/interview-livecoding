package private

import (
	"github.com/farawaygg/wisdom/pkg/wisdom"
)

var _ wisdom.WisdomSvcServer = (*Service)(nil)

type Service struct {
	wisdom.UnimplementedWisdomSvcServer

	wisdoms Wisdom
}

func New(w Wisdom) *Service {
	return &Service{
		wisdoms: w,
	}
}
