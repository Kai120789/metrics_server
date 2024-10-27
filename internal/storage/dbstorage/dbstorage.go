package dbstorage

import (
	"context"
	"fmt"
	"server/internal/dto"
	"server/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v5"
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

	var retMetrics []models.Metric
	next := s.getNextPollCountDelta()

	if next != nil {
		*metrics[0].Delta = *next
	}

	fmt.Println(metrics[0].Delta)

	for _, metric := range metrics {
		var retMetric models.Metric
		query := `INSERT INTO metrics (name, type, value, delta) VALUES ($1, $2, $3, $4) RETURNING id, name, type, value, delta, created_at`
		err := s.Conn.QueryRow(context.Background(), query, metric.Name, metric.Type, metric.Value, &metric.Delta).Scan(&retMetric.ID, &retMetric.Name, &retMetric.Type, &retMetric.Value, &retMetric.Delta, &retMetric.CreatedAt)
		if err != nil {
			s.Logger.Error("Failed to insert metric", zap.Error(err))
			continue
		}

		retMetrics = append(retMetrics, retMetric)
	}

	return &retMetrics, nil
}

func (s *Storage) SetMetric(metric dto.Metric) (*models.Metric, error) {
	var retMetric models.Metric

	query := `INSERT INTO metrics (name, type, value, delta) VALUES ($1, $2, $3, $4) RETURNING id, name, type, value, delta, created_at`
	err := s.Conn.QueryRow(context.Background(), query, metric.Name, metric.Type, metric.Value, nil).Scan(&retMetric.ID, &retMetric.Name, &retMetric.Type, &retMetric.Value, &retMetric.Delta, &retMetric.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &retMetric, nil
}

func (s *Storage) GetMetricValue(name string, typeStr string) (*int64, error) {
	var value int64

	query := `SELECT value FROM metrics WHERE name = $1 AND type = $2 ORDER BY created_at DESC`
	row := s.Conn.QueryRow(context.Background(), query, name, typeStr)

	err := row.Scan(&value)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &value, nil
}

func (s *Storage) GetMetricsForHTML() (*[]models.Metric, error) {
	query := `SELECT * FROM metrics ORDER BY created_at DESC LIMIT 31`
	rows, err := s.Conn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var metrics []models.Metric
	for rows.Next() {
		var metric models.Metric
		err := rows.Scan(&metric.ID, &metric.Name, &metric.Type, &metric.Value, &metric.Delta, &metric.CreatedAt)
		if err != nil {
			return nil, err
		}

		metrics = append(metrics, metric)
	}

	return &metrics, nil
}

func (s *Storage) getNextPollCountDelta() *int64 {
	query := `SELECT delta FROM metrics WHERE type = 'counter' ORDER BY delta DESC LIMIT 1`
	row := s.Conn.QueryRow(context.Background(), query)

	var max int64
	err := row.Scan(&max)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	query = `SELECT delta FROM metrics WHERE type = 'counter' LIMIT 1`
	row = s.Conn.QueryRow(context.Background(), query)

	var min int64
	err = row.Scan(&min)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	next := max + min

	fmt.Println(min, max, next)

	return &next
}
