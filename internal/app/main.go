package app

import (
	"fmt"
	"server/internal/config"
	"server/internal/storage/dbstorage"
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

	// connect to db postgres
	dbConn, err := dbstorage.Connection(cfg.DBDSN)
	if err != nil {
		log.Fatal("error connect to db", zap.Error(err))
	}

	defer dbConn.Close()

	// init storage

	// init service

	// init service

	// init handler

	// init router

	// start http-server
}
