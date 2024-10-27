package service

import (
	"server/internal/dto"
	"server/internal/models"
	"text/template"
)

type Service struct {
	storage Storager
}

type Storager interface {
	SetUpdates(metrics []dto.Metric) (*[]models.Metric, error)
	SetMetric(metric dto.Metric) (*models.Metric, error)
	GetMetricValue(name string, typeStr string) (*int64, error)
	GetMetricsForHTML() (*[]models.Metric, error)
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

func (s *Service) GetMetricValue(name string, typeStr string) (*int64, error) {
	val, err := s.storage.GetMetricValue(name, typeStr)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (s *Service) GetHTML() (*[]models.Metric, *template.Template, error) {
	metrics, err := s.storage.GetMetricsForHTML()
	if err != nil {
		return nil, nil, err
	}

	// Подготовка шаблона HTML для отображения метрик
	tmpl := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Metrics</title>
	</head>
	<body>
		<h1>Metrics</h1>
		<table border="1">
			<tr>
				<th>ID</th>
				<th>Name</th>
				<th>Type</th>
				<th>Value</th>
				<th>Delta</th>
				<th>Created At</th>
			</tr>
			{{range .}}
			<tr>
				<td>{{.ID}}</td>
				<td>{{.Name}}</td>
				<td>{{.Type}}</td>
				<td>{{.Value}}</td>
				<td>{{.Delta}}</td>
				<td>{{.CreatedAt}}</td>
			</tr>
			{{end}}
		</table>
	</body>
	</html>`

	t, err := template.New("metrics").Parse(tmpl)
	if err != nil {
		return nil, nil, err
	}

	return metrics, t, nil
}
