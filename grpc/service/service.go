package service

import (
	"context"
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

func (f *OrderService) Create(ctx context.Context, req *order_service.CreateOrder) (*order_service.Order, error) {
	f.log.Info("Create Order: ", logger.Any("req", req))

	resp, err := f.strg.Order().Create(ctx, req)

	if err != nil {
		f.log.Error("Create Order: ", logger.Error(err))
		return &order_service.Order{}, err
	}
	return resp, nil
}

func (f *OrderService) GetById(ctx context.Context, req *order_service.OrderPrimaryKey) (*order_service.Order, error) {
	f.log.Info("Get Single Order: ", logger.Any("req", req))

	resp, err := f.strg.Order().GetById(ctx, req)

	if err != nil {
		f.log.Error("failed to get single order: ", logger.Error(err))
		return &order_service.Order{}, err
	}
	return resp, nil
}