package app

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"server/internal/config"
	"server/internal/service"
	"server/internal/storage"
	"server/internal/storage/dbstorage"
	"server/internal/transport/grpc/proto"
	"server/internal/transport/grpc/proto/server"
	"server/internal/transport/http/handler"
	"server/internal/transport/http/router"
	"server/internal/utils"
	"server/pkg/logger"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
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
	serv := service.New(dbstor, log)

	// init handler
	handl := handler.New(serv, log, cfg)

	// init router
	r := router.New(&handl)

	// start http-server
	log.Info("starting server", zap.String("address", cfg.ServerAddress))

	go func() {
		grpcServer := grpc.NewServer()

		// create and reg grpc-server
		grpcServerInstance := server.NewGRPCServer(serv)
		proto.RegisterMetricServiceServer(grpcServer, grpcServerInstance)

		listener, err := net.Listen("tcp", ":50051")
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println("server gRPC start on :50051")
		if err := grpcServer.Serve(listener); err != nil {
			fmt.Println(err.Error())
		}
	}()

	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: r,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", zap.Error(err))
	}

}
