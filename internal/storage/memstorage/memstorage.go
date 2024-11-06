package memstorage

import (
	"fmt"
	"server/internal/dto"
	"server/internal/models"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Storage struct {
	Metrics []models.Metric
	Logger  *zap.Logger
	mu      sync.Mutex
}

func New(metrics []models.Metric, log *zap.Logger) *Storage {
	return &Storage{
		Metrics: metrics,
		Logger:  log,
	}
}

func (s *Storage) SetUpdates(metrics []dto.Metric) ([]models.Metric, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var id uint = 1
	if len(s.Metrics) != 0 {
		if s.Metrics[0].Delta != nil {
			*(s.Metrics)[0].Delta += 5
		}

		for i := range s.Metrics {
			s.Metrics[i].Value = (metrics)[i].Value
			s.Metrics[i].CreatedAt = time.Now()
		}
	} else {
		for _, metric := range metrics {

			var retMetric models.Metric = models.Metric{
				ID:        id,
				Name:      metric.Name,
				Type:      metric.Type,
				Value:     metric.Value,
				Delta:     metric.Delta,
				CreatedAt: time.Now(),
			}

			id += 1

			s.Metrics = append(s.Metrics, retMetric)
		}
	}

	fmt.Println(*(s.Metrics)[0].Delta)

	return s.Metrics, nil
}

func (s *Storage) SetMetric(metric dto.Metric) (*models.Metric, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var retMetric models.Metric
	var flag bool = true

	for _, met := range s.Metrics {
		if met.Name == metric.Name && met.Type == metric.Type {
			met.Value = metric.Value
			met.Delta = metric.Delta
			retMetric = met
			flag = false
		}
	}

	if flag {
		retMetric = models.Metric{
			ID:        uint(len(s.Metrics) + 1),
			Name:      metric.Name,
			Type:      metric.Type,
			Value:     metric.Value,
			Delta:     metric.Delta,
			CreatedAt: time.Now(),
		}
	}

	s.Metrics = append(s.Metrics, retMetric)

	return &retMetric, nil
}

func (s *Storage) GetMetricValue(name string, typeStr string) (*float64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var value float64

	for _, metric := range s.Metrics {
		if metric.Name == name && typeStr == metric.Type {
			value = float64(*metric.Value)
		}
	}

	return &value, nil
}

func (s *Storage) GetMetricsForHTML() ([]models.Metric, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Metrics, nil
}
