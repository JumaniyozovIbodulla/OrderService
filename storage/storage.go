package storage

import (
	"context"
	"order/genproto/order_service"
	orp "order/genproto/order_product_service"
	orn "order/genproto/order_notes"

)

type IStorage interface {
	CloseDB()
	Order() OrderRepo
	OrderProduct() OrderProductsRepo
	OrderNotes() OrderNotesRepo
}

type OrderRepo interface {
	Create(ctx context.Context, req *order_service.CreateOrder) (resp *order_service.Order, err error)
	GetById(ctx context.Context, req *order_service.OrderPrimaryKey) (resp *order_service.Order, err error)
	GetAll(ctx context.Context, req *order_service.GetListOrderRequest) (resp *order_service.GetListOrderResponse, err error)
	Update(ctx context.Context, req *order_service.UpdateOrder) (resp *order_service.Order, err error)
	Delete(ctx context.Context, req *order_service.OrderPrimaryKey) (resp *order_service.Empty, err error)
}

type OrderProductsRepo interface {
	Create(ctx context.Context, req *orp.CreateOrderProduct) (resp *orp.OrderProduct, err error)
	// GetById(ctx context.Context, req *orp.OrderProductPrimaryKey) (*orp.OrderProduct, error)
	// GetAll(ctx context.Context, req *orp.GetListOrderProductRequest) (resp *orp.GetListOrderProductResponse, err error)
	// Update(ctx context.Context, req *orp.UpdateOrderProduct) (resp *orp.OrderProduct, err error)
	// Delete(ctx context.Context, req *orp.OrderProductPrimaryKey) (resp *orp.Empty, err error)
}

type OrderNotesRepo interface {
	Create(ctx context.Context, req *orn.CreateOrderNotes) (resp *orn.OrderNotes, err error)
	// GetById(ctx context.Context, req *orn.OrderNotesPrimaryKey) (*orn.OrderNotes, error)
	// GetAll(ctx context.Context, req *orn.GetListOrderNotesRequest) (resp *orn.GetListOrderNotesResponse, err error)
	// Update(ctx context.Context, req *orn.UpdateOrderNotes) (resp *orn.OrderNotes, err error)
	// Delete(ctx context.Context, req *orn.OrderNotesPrimaryKey) (resp *orn.Empty, err error)
}