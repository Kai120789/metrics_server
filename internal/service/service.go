package service

import (
	"server/internal/dto"
	"server/internal/models"
)

type Service struct {
	storage Storager
}

type Storager interface {
	SetUpdates(metrics []dto.Metric) (*[]models.Metric, error)
	SetMetric(metric dto.Metric) (*models.Metric, error)
	GetMetricValue()
	GetHTML()
}

func New(s Storager) *Service {
	return &Service{
		storage: s,
	}
}

func (s *Service) SetUpdates(metrics []dto.Metric) (*[]models.Metric, error) {
	met, err := s.storage.SetUpdates(metrics)
	if err != nil {
		return nil, err
	}

	return met, nil
}

func (s *Service) SetMetric(metric dto.Metric) (*models.Metric, error) {
	met, err := s.storage.SetMetric(metric)
	if err != nil {
		return nil, err
	}

	return met, nil
}

func (s *Service) GetMetricValue() {

}

func (s *Service) GetHTML() {

}
