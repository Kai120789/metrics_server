package service

import (
	"fmt"
	"net/http"
	"server/internal/dto"
	"server/internal/models"
	"text/template"
)

type Service struct {
	storage Storager
}

type Storager interface {
	SetUpdates(metrics []dto.Metric) ([]models.Metric, error)
	SetMetric(metric dto.Metric) (*models.Metric, error)
	GetMetricValue(name string, typeStr string) (*float64, error)
	GetMetricsForHTML() ([]models.Metric, error)
}

func New(s Storager) *Service {
	return &Service{
		storage: s,
	}
}

func (s *Service) SetUpdates(metrics []dto.Metric) ([]models.Metric, error) {
	met, err := s.storage.SetUpdates(metrics)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return met, nil
}

func (s *Service) SetMetric(metric dto.Metric) (*models.Metric, error) {
	met, err := s.storage.SetMetric(metric)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return met, nil
}

func (s *Service) GetMetricValue(name string, typeStr string) (*float64, error) {
	val, err := s.storage.GetMetricValue(name, typeStr)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return val, nil
}

func (s *Service) GetHTML(w http.ResponseWriter) error {
	metrics, err := s.storage.GetMetricsForHTML()
	if err != nil {
		fmt.Println(err.Error())
		return err
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
		fmt.Println(err.Error())
		return err
	}

	err = t.Execute(w, metrics)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}
