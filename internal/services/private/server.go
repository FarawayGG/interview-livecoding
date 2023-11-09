package private

import (
	"github.com/farawaygg/wisdom/pkg/wisdom"
)

var _ wisdom.WisdomSvcServer = (*Service)(nil)

type Service struct {
	wisdoms Wisdom
}

func New(w Wisdom) *Service {
	return &Service{
		wisdoms: w,
	}
}
