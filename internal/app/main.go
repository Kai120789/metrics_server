package app

import (
	"fmt"
	"net/http"
	"server/internal/config"
	"server/internal/service"
	"server/internal/storage"
	"server/internal/transport/http/handler"
	"server/internal/transport/http/router"
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

	/*_, err = filestorage.CreateFile(cfg.FilePath)
	if err != nil {
		return
	}*/

	// connect to db postgres
	/*dbConn, err := dbstorage.Connection(cfg.DBDSN)
	if err != nil {
		log.Fatal("error connect to db", zap.Error(err))
	}

	defer dbConn.Close()*/

	// init storage
	dbstor := storage.New(nil, log, cfg)

	// init service
	serv := service.New(dbstor)

	// init handler
	handl := handler.New(serv, log, cfg)

	// init router
	r := router.New(&handl)

	// start http-server
	log.Info("starting server", zap.String("address", "localhost:8082"))

	srv := &http.Server{
		Addr:    "localhost:8082",
		Handler: r,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", zap.Error(err))
	}
}
