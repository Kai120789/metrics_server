package service

import (
	"fmt"
	"server/internal/dto"
	"server/internal/models"
)

type Service struct {
	storage Storager
}

type Storager interface {
	SetUpdates(metrics []dto.Metric) (*[]models.Metric, error)
	SetUpdate()
	SetMetric()
	GetMetricValue()
	GetHTML()
}

func New(s Storager) *Service {
	return &Service{
		storage: s,
	}
}

func (s *Service) SetUpdates(metrics []dto.Metric) (*[]models.Metric, error) {
	fmt.Println(2)
	met, err := s.storage.SetUpdates(metrics)
	if err != nil {
		return nil, err
	}

	return met, nil
}

func (s *Service) SetUpdate() {

}

func (s *Service) SetMetric() {

}

func (s *Service) GetMetricValue() {

}

func (s *Service) GetHTML() {

}
