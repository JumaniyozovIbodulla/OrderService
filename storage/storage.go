package storage

import (
	"context"
	"order/genproto/order_service"
)

type IStorage interface {
	CloseDB()
	Order() OrderRepo
}

type OrderRepo interface {
	Create(ctx context.Context, req *order_service.CreateOrder) (*order_service.Order, error)
	GetById(ctx context.Context, req *order_service.OrderPrimaryKey) (*order_service.Order, error)
	// GetAll(ctx context.Context, req *order_service.GetListOrderRequest) (resp *order_service.GetListOrderResponse, err error)
	// Update(ctx context.Context, req *order_service.UpdateOrder) (resp *order_service.Order, err error)
	// Delete(ctx context.Context, req *order_service.OrderPrimaryKey) (err error)
}