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
	Create(ctx context.Context, req *order_service.CreateOrder) (*order_service.Order, error)
	GetById(ctx context.Context, req *order_service.OrderPrimaryKey) (*order_service.Order, error)
	GetAll(ctx context.Context, req *order_service.GetListOrderRequest) (*order_service.GetListOrderResponse, error)
	Update(ctx context.Context, req *order_service.UpdateOrder) (*order_service.Order, error)
	Delete(ctx context.Context, req *order_service.OrderPrimaryKey) (*order_service.Empty, error)
}

type OrderProductsRepo interface {
	Create(ctx context.Context, req *orp.CreateOrderProduct) (*orp.OrderProduct, error)
	GetById(ctx context.Context, req *orp.OrderProductPrimaryKey) (*orp.OrderProduct, error)
	GetAll(ctx context.Context, req *orp.GetListOrderProductRequest) (*orp.GetListOrderProductResponse, error)
	Update(ctx context.Context, req *orp.UpdateOrderProduct) (*orp.OrderProduct, error)
	Delete(ctx context.Context, req *orp.OrderProductPrimaryKey) (*orp.Empty, error)
}

type OrderNotesRepo interface {
	Create(ctx context.Context, req *orn.CreateOrderNotes) (*orn.OrderNotes, error)
	GetById(ctx context.Context, req *orn.OrderNotesPrimaryKey) (*orn.OrderNotes, error)
	GetAll(ctx context.Context, req *orn.GetListOrderNotesRequest) (*orn.GetListOrderNotesResponse, error)
	Update(ctx context.Context, req *orn.UpdateOrderNotes) (*orn.OrderNotes, error)
	Delete(ctx context.Context, req *orn.OrderNotesPrimaryKey) (*orn.Empty, error)
}