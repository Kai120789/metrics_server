package dbstorage

import (
	"context"
	"fmt"
	"server/internal/dto"
	"server/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type Storage struct {
	Conn   *pgxpool.Pool
	Logger *zap.Logger
}

func New(dbConn *pgxpool.Pool, log *zap.Logger) *Storage {
	return &Storage{
		Conn:   dbConn,
		Logger: log,
	}
}

func Connection(connectionStr string) (*pgxpool.Pool, error) {
	db, err := pgxpool.Connect(context.Background(), connectionStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to db: %v", err)
	}

	return db, nil
}

func (s *Storage) SetUpdates(metrics []dto.Metric) (*[]models.Metric, error) {
	fmt.Println(3)
	var retMetrics []models.Metric

	for _, metric := range metrics {
		var retMetric models.Metric
		query := `INSERT INTO metrics (name, type, value, delta) VALUES ($1, $2, $3, $4) RETURNING id, name, type, value, delta, created_at`
		err := s.Conn.QueryRow(context.Background(), query, metric.Name, metric.Type, &metric.Value, &metric.Delta).Scan(&retMetric.ID, &retMetric.Name, &retMetric.Type, &retMetric.Value, &retMetric.Delta, &retMetric.CreatedAt)
		if err != nil {
			s.Logger.Error("Failed to insert metric", zap.Error(err))
			continue // Переход к следующей метрике
		}

		retMetrics = append(retMetrics, retMetric)
	}

	return &retMetrics, nil
}

func (s *Storage) SetUpdate() {

}

func (s *Storage) SetMetric() {

}

func (s *Storage) GetMetricValue() {

}

func (s *Storage) GetHTML() {

}
