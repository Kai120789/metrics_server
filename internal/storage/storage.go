package storage

import (
	"server/internal/config"
	"server/internal/dto"
	"server/internal/models"
	"server/internal/storage/dbstorage"
	"server/internal/storage/filestorage"
	"server/internal/storage/memstorage"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type Storage interface {
	SetUpdates(metrics []dto.Metric) ([]models.Metric, error)
	SetMetric(metric dto.Metric) (*models.Metric, error)
	GetMetricValue(name string, typeStr string) (*float64, error)
	GetMetricsForHTML() ([]models.Metric, error)
}

func New(dbConn *pgxpool.Pool, log *zap.Logger, cfg *config.Config) Storage {
	switch {
	case dbConn != nil:
		log.Info("used postgres")
		return dbstorage.New(dbConn, log)
	case cfg.FilePath != "":
		log.Info("used file")
		return filestorage.New(cfg.FilePath, log)
	default:
		log.Info("used memory")
		return memstorage.New([]models.Metric{}, &zap.Logger{})
	}
}
