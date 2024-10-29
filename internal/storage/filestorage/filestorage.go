package filestorage

import (
	"encoding/json"
	"fmt"
	"os"
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

func CreateFile(filePath string) (*os.File, error) {
	var file os.File
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		defer file.Close()
		fmt.Println("File successfuly created")

	} else {
		fmt.Println("File is ready exist")
	}

	return &file, nil
}

func (s *Storage) SetUpdates(metrics []dto.Metric) (*[]models.Metric, error) {
	retMetricsPrev, err := s.readMetrics()
	if err != nil {
		return nil, err
	}

	delta := retMetricsPrev[0].Delta
	*metrics[0].Delta = *delta + 5

	data, err := json.Marshal(metrics)
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
	return nil, nil
}

func (s *Storage) GetMetricValue(name string, typeStr string) (*int64, error) {
	return nil, nil
}

func (s *Storage) GetMetricsForHTML() (*[]models.Metric, error) {
	return nil, nil
}

func (s *Storage) readMetrics() ([]models.Metric, error) {
	file, err := os.ReadFile(s.FilePath)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	// Проверка на пустой файл
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
