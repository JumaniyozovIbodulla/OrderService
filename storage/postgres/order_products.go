package postgres

import (
	"context"
	orp "order/genproto/order_product_service"
	"order/storage"

	"github.com/jackc/pgx/v4/pgxpool"
)

type orderProducsRepo struct {
	db *pgxpool.Pool
}

func NewOrderProductsRepo(db *pgxpool.Pool) storage.OrderProductsRepo {
	return &orderProducsRepo{
		db: db,
	}
}

func (o *orderProducsRepo) Create(ctx context.Context, req *orp.CreateOrderProduct) (resp *orp.OrderProduct, err error) {
	return
}
