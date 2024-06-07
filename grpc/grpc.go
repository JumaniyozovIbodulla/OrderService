package grpc

import (
	"order/config"
	"order/genproto/order_notes"
	"order/genproto/order_product_service"
	"order/genproto/order_service"
	"order/grpc/client"
	"order/grpc/service"
	"order/storage"

	"github.com/saidamir98/udevs_pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func SetUpServer(cfg config.Config, log logger.LoggerI, strg storage.IStorage, ordSerManager client.OrderServiceManager, ordProductManager client.OrderProductsManager, ordNotes client.OrderNotesManager) *grpc.Server {
	grpcServer := grpc.NewServer()

	order_service.RegisterOrderServiceServer(grpcServer, service.NewOrderService(cfg, log, strg, ordSerManager))
	order_product_service.RegisterOrderProductsServer(grpcServer, service.NewOrderProductService(cfg, log, strg, ordProductManager))
	order_notes.RegisterOrderStatusNotesServer(grpcServer, service.NewOrderNotesService(cfg, log, strg, ordNotes))
	reflection.Register(grpcServer)
	return grpcServer
}
