package main

import (
	"context"
	"net"
	"order/config"
	"order/grpc"
	"order/grpc/client"
	"order/storage/postgres"

	"github.com/saidamir98/udevs_pkg/logger"
)

func main() {
	cfg := config.Load()

	loggerLevel := logger.LevelDebug

	switch cfg.Environment {
	case config.DebugMode:
		loggerLevel = logger.LevelDebug
	case config.TestMode:
		loggerLevel = logger.LevelDebug
	default:
		loggerLevel = logger.LevelInfo
	}

	log := logger.NewLogger(cfg.ServiceName, loggerLevel)

	defer logger.Cleanup(log)

	pgStore, err := postgres.NewPostgres(context.Background(), cfg)

	if err != nil {
		log.Panic("postgres.NewPostgres", logger.Error(err))
	}

	defer pgStore.CloseDB()

	order, err := client.NewGrpcOrderClients(cfg)

	if err != nil {
		log.Panic("client.NewGrpcClients: ", logger.Error(err))
	}

	orderProduct, err := client.NewGrpcOrderProductClients(cfg)

	if err != nil {
		log.Panic("client.NewGrpcOrderProductClients: ", logger.Error(err))
	}

	orderNotes, err := client.NewGrpcOrderNotesClients(cfg)

	if err != nil {
		log.Panic("client.NewGrpcOrderNotesClients: ", logger.Error(err))
	}

	grpcServer := grpc.SetUpServer(cfg, log, pgStore, order, orderProduct, orderNotes)

	lis, err := net.Listen("tcp", cfg.ContentGRPCPort)

	if err != nil {
		log.Panic("net.Listen: ", logger.Error(err))
	}

	log.Info("GRPC: Server being started...", logger.String("port", cfg.ContentGRPCPort))

	if err := grpcServer.Serve(lis); err != nil {
		log.Panic("grpcServer.Serve: ", logger.Error(err))
	}
}