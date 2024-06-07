package service

import (
	"order/config"
	"order/genproto/order_notes"
	"order/genproto/order_product_service"
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

type OrderProductService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.IStorage
	services client.OrderProductsManager
	*order_product_service.UnimplementedOrderProductsServer
}

type OrderNotesService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.IStorage
	services client.OrderNotesManager
	*order_notes.UnimplementedOrderStatusNotesServer
}

func NewOrderService(cfg config.Config, log logger.LoggerI, strg storage.IStorage, srvc client.OrderServiceManager) *OrderService {
	return &OrderService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvc,
	}
}

func NewOrderProductService(cfg config.Config, log logger.LoggerI, strg storage.IStorage, srvc client.OrderProductsManager) *OrderProductService {
	return &OrderProductService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvc,
	}
}

func NewOrderNotesService(cfg config.Config, log logger.LoggerI, strg storage.IStorage, srvc client.OrderNotesManager) *OrderNotesService {
	return &OrderNotesService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvc,
	}
}
