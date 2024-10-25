package service

type Service struct {
	storage Storager
}

type Storager interface {
}

func New(s Storager) *Service {
	return &Service{
		storage: s,
	}
}
