package filestorage

import (
	"encoding/json"
	"fmt"
	"os"
	"server/internal/dto"
	"server/internal/models"
	"time"

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

func CreateFile(filePath string) (*os.File, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		defer file.Close()
		fmt.Println("File successfuly created")
		return file, nil

	} else {
		fmt.Println("File is ready exist")
	}

	return nil, nil
}

func (s *Storage) SetUpdates(metrics []dto.Metric) (*[]models.Metric, error) {
	retMetricsPrev, err := s.readMetrics()
	if err != nil {
		return nil, err
	}

	if len(retMetricsPrev) != 0 {
		delta := retMetricsPrev[0].Delta
		*metrics[0].Delta = *delta + 5
	}

	var dataMetrics []models.Metric
	var id uint = 1

	for _, metric := range metrics {

		dataMetric := models.Metric{
			ID:        id,
			Name:      metric.Name,
			Type:      metric.Type,
			Value:     metric.Value,
			Delta:     metric.Delta,
			CreatedAt: time.Now(),
		}

		dataMetrics = append(dataMetrics, dataMetric)
		id += 1
	}

	data, err := json.Marshal(dataMetrics)
	if err != nil {
		return nil, err
	}

	os.WriteFile(s.FilePath, data, os.ModePerm)

	retMetrics, err := s.readMetrics()
	if err != nil {
		return nil, err
	}

	return &retMetrics, nil
}

func (s *Storage) SetMetric(metric dto.Metric) (*models.Metric, error) {
	retMetrics, err := s.readMetrics()
	if err != nil {
		return nil, err
	}

	var flag bool = true

	var retMetric models.Metric

	for _, met := range retMetrics {
		if met.Name == metric.Name && met.Type == metric.Type {
			met.Value = metric.Value
			met.Delta = metric.Delta
			retMetric = met
			flag = false
		}
	}

	if flag {
		retMetric = models.Metric{
			ID:        uint(len(retMetrics) + 1),
			Name:      metric.Name,
			Type:      metric.Type,
			Value:     metric.Value,
			Delta:     metric.Delta,
			CreatedAt: time.Now(),
		}
	}

	retMetrics = append(retMetrics, retMetric)

	data, err := json.Marshal(retMetrics)
	if err != nil {
		return nil, err
	}

	os.WriteFile(s.FilePath, data, os.ModePerm)

	return &retMetric, nil
}

func (s *Storage) GetMetricValue(name string, typeStr string) (*int64, error) {
	retMetrics, err := s.readMetrics()
	if err != nil {
		return nil, err
	}

	var value int64

	for _, metric := range retMetrics {
		if metric.Name == name && typeStr == metric.Type {
			value = int64(*metric.Value)
		}
	}

	return &value, nil
}

func (s *Storage) GetMetricsForHTML() (*[]models.Metric, error) {
	retMetrics, err := s.readMetrics()
	if err != nil {
		return nil, err
	}

	return &retMetrics, nil
}

func (s *Storage) readMetrics() ([]models.Metric, error) {
	file, err := os.ReadFile(s.FilePath)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	if len(file) == 0 {
		return []models.Metric{}, nil
	}

	var metrics []models.Metric
	err = json.Unmarshal(file, &metrics)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return metrics, nil
}
