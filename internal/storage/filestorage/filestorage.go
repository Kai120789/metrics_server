package filestorage

import (
	"server/internal/dto"
	"server/internal/models"

	"go.uber.org/zap"
)

type Storage struct {
	FilePath string
	Logger   *zap.Logger
}

func New(fp string, log *zap.Logger) *Storage {
	return &Storage{
		FilePath: fp,
		Logger:   log,
	}
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
