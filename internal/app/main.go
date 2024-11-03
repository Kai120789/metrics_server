package app

import (
	"fmt"
	"net/http"
	"server/internal/config"
	"server/internal/service"
	"server/internal/storage"
	"server/internal/storage/dbstorage"
	"server/internal/storage/filestorage"
	"server/internal/transport/http/handler"
	"server/internal/transport/http/router"
	"server/internal/utils"
	"server/pkg/logger"

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
		_, err = filestorage.CreateFile(cfg.FilePath)
		if err != nil {
			return
		}
	}

	// connect to db postgres
	if cfg.DBDSN != "" {
		// flag for migrations
		if cfg.Migrations {
			utils.DoMigrate()
		}

		dbConn, err := dbstorage.Connection(cfg.DBDSN)
		if err != nil {
			log.Fatal("error connect to db", zap.Error(err))
		}

		defer dbConn.Close()
	}

	// init storage
	dbstor := storage.New(nil, log, cfg)

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
