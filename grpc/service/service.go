package service

import (
	"order/config"
	"order/genproto/order_service"
	"order/grpc/client"
	"order/storage"

	"github.com/saidamir98/udevs_pkg/logger"
)

type OrderService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.IStorage
	services client.OrderServiceManager
	*order_service.UnimplementedOrderServiceServer
}

func NewOrderService(cfg config.Config, log logger.LoggerI, strg storage.IStorage, srvc client.OrderServiceManager) *OrderService {
	return &OrderService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvc,
	}
}