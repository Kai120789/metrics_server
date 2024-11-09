package main

import (
	"fmt"
	"net"
	"server/internal/config"
	"server/internal/service"
	"server/internal/storage"
	"server/internal/storage/dbstorage"
	"server/internal/transport/grpc/proto"
	"server/internal/transport/grpc/proto/server"
	"server/internal/utils"
	"server/pkg/logger"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
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

	var dbConn *pgxpool.Pool

	// connect to db postgres
	if cfg.DBDSN != "" {
		// flag for migrations
		if cfg.Migrations {
			utils.DoMigrate()
		}

		dbConn, err = dbstorage.Connection(string(cfg.DBDSN[:27]) + "localhost:5433/metrics")
		if err != nil {
			log.Fatal("error connect to db", zap.Error(err))
		}

		defer dbConn.Close()
	}

	// Создание экземпляра Service с нужным хранилищем
	storage := storage.New(dbConn, log, cfg)
	svc := service.New(storage)
	grpcServer := grpc.NewServer()

	// Создание и регистрация gRPC-сервера
	grpcServerInstance := server.NewGRPCServer(svc)
	proto.RegisterMetricServiceServer(grpcServer, grpcServerInstance)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("Сервер gRPC запущен на :50051")
	if err := grpcServer.Serve(listener); err != nil {
		fmt.Println(err.Error())
	}
}
