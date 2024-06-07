package service

import (
	"context"
	"order/genproto/order_service"

	"github.com/saidamir98/udevs_pkg/logger"
)

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

func (f *OrderService) Update(ctx context.Context, req *order_service.UpdateOrder) (*order_service.Order, error) {
	f.log.Info("Update an Order: ", logger.Any("req", req))

	resp, err := f.strg.Order().Update(ctx, req)

	if err != nil {
		f.log.Error("Update an Order: ", logger.Error(err))
		return &order_service.Order{}, err
	}
	return resp, nil
}

func (o *OrderService) Delete(ctx context.Context, req *order_service.OrderPrimaryKey) (*order_service.Empty, error) {
	o.log.Info("Delete an order: ", logger.Any("req", req))

	resp, err := o.strg.Order().Delete(ctx, req)

	if err != nil {
		o.log.Error("Delete an Order: ", logger.Error(err))
		return nil, err
	}
	return resp, nil
}

func (f *OrderService) GetAll(ctx context.Context, req *order_service.GetListOrderRequest) (*order_service.GetListOrderResponse, error) {
	f.log.Info("Get All Orders: ", logger.Any("req", req))

	resp, err := f.strg.Order().GetAll(ctx, req)

	if err != nil {
		f.log.Error("failed to get all orders: ", logger.Error(err))
		return nil, err
	}
	return resp, nil
}
