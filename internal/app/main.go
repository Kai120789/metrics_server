package app

import (
	"fmt"
	"net/http"
	"os"
	"server/internal/config"
	"server/internal/service"
	"server/internal/storage"
	"server/internal/storage/dbstorage"
	"server/internal/transport/http/handler"
	"server/internal/transport/http/router"
	"server/internal/utils"
	"server/pkg/logger"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

func StartServer() {
	// init config
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(cfg)

	// init logger
	zapLog, err := logger.New(cfg.LogLevel)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	log := zapLog.ZapLogger

	if cfg.FilePath != "" {
		// check is file exist
		if _, err := os.Stat(cfg.FilePath); os.IsNotExist(err) {
			file, err := os.Create(cfg.FilePath)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			defer file.Close()
			fmt.Println("file succesfully created:", cfg.FilePath)

		} else {
			fmt.Println("file is exist now:", cfg.FilePath)
		}
	}

	var dbConn *pgxpool.Pool

	// connect to db postgres
	if cfg.DBDSN != "" {
		// flag for migrations
		if cfg.Migrations {
			utils.DoMigrate()
		}

		dbConn, err = dbstorage.Connection(cfg.DBDSN)
		if err != nil {
			log.Fatal("error connect to db", zap.Error(err))
		}

		defer dbConn.Close()
	}

	// init storage
	dbstor := storage.New(dbConn, log, cfg)

	if cfg.RestoreMetrics {
		allMetrics, err := dbstor.GetMetricsForHTML()
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println(allMetrics)
	}

	// init service
	serv := service.New(dbstor)

	// init handler
	handl := handler.New(serv, log, cfg)

	// init router
	r := router.New(&handl)

	// start http-server
	log.Info("starting server", zap.String("address", cfg.ServerAddress))

	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: r,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", zap.Error(err))
	}
}
