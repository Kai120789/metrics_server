package app

import (
	"fmt"
	"server/internal/config"
	"server/pkg/logger"
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

	_ = log

	// connect to db postgres

	// init storage

	// init service

	// init service

	// init handler

	// init router

	// start http-server
}
