package storage

import (
	"server/internal/config"
	"server/internal/storage/dbstorage"
	"server/internal/storage/filestorage"
	"server/internal/storage/memstorage"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type Storage interface {
}

func New(dbConn *pgxpool.Pool, log *zap.Logger, cfg *config.Config, value string) Storage {
	if value == cfg.DBDSN {
		return dbstorage.New(dbConn, log)
	} else if value == cfg.FilePath {
		return filestorage.New(cfg.FilePath, log)
	}

	return memstorage.New()
}
