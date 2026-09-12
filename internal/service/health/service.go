package health

import "context"

type Pinger interface {
	Ping(ctx context.Context) error
}

type Service struct {
	isReady bool
	pinger  Pinger
}

func NewService(pinger Pinger) *Service {
	return &Service{pinger: pinger}
}

func (s *Service) IsReady(ctx context.Context) bool {
	if err := s.pinger.Ping(ctx); err != nil {
		return false
	}
	if !s.isReady {
		return false
	}
	return true
}

func (s *Service) SetReady() {
	s.isReady = true
}
