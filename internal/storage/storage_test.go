package storage_test

import (
	"server/internal/config"
	"server/internal/storage"
	"server/internal/storage/dbstorage"
	"server/internal/storage/filestorage"
	"server/internal/storage/memstorage"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

func TestNewStorage(t *testing.T) {
	log, _ := zap.NewProduction()

	t.Run("if dbConn != nil", func(t *testing.T) {
		mockDbConn := &pgxpool.Pool{}
		cfg := &config.Config{}

		storage := storage.New(mockDbConn, log, cfg)
		if _, ok := storage.(*dbstorage.Storage); !ok {
			t.Errorf("expected dbstorage.DBStorage, got %T", storage)
		}
	})

	t.Run("if filePath != nil", func(t *testing.T) {
		mockFilePath := "mockPath.go"
		cfg := &config.Config{}
		cfg.FilePath = mockFilePath

		storage := storage.New(nil, log, cfg)
		if _, ok := storage.(*filestorage.Storage); !ok {
			t.Errorf("expected filestorage.FileStorage, got %T", storage)
		}
	})

	t.Run("default", func(t *testing.T) {
		cfg := &config.Config{}

		storage := storage.New(nil, log, cfg)
		if _, ok := storage.(*memstorage.Storage); !ok {
			t.Errorf("expected memstorage.MemStorage, got %T", storage)
		}
	})
}
