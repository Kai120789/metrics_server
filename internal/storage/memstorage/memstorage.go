package memstorage

import (
	"server/internal/dto"
	"server/internal/models"
)

type Storage struct {
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) SetUpdates(metrics []dto.Metric) (*[]models.Metric, error) {
	return nil, nil
}

func (s *Storage) SetMetric(metric dto.Metric) (*models.Metric, error) {
	return nil, nil
}

func (s *Storage) GetMetricValue() {

}

func (s *Storage) GetHTML() {

}
