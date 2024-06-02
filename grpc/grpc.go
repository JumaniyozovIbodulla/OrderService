package grpc

import (
	"order/config"
	"order/genproto/order_service"
	"order/grpc/client"
	"order/grpc/service"
	"order/storage"

	"github.com/saidamir98/udevs_pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func SetUpServer(cfg config.Config, log logger.LoggerI, strg storage.IStorage, srvc client.OrderServiceManager) (grpcServer *grpc.Server) {
	grpcServer = grpc.NewServer()

	order_service.RegisterOrderServiceServer(grpcServer, service.NewOrderService(cfg, log, strg, srvc))
	reflection.Register(grpcServer)
	return 
}
