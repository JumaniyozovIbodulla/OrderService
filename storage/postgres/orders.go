package postgres

import (
	"context"
	"log"
	"order/genproto/order_service"
	"order/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type orderRepo struct {
	db *pgxpool.Pool
}

func NewOrderRepo(db *pgxpool.Pool) storage.OrderRepo {
	return &orderRepo{
		db: db,
	}
}

func (o *orderRepo) Create(ctx context.Context, req *order_service.CreateOrder) (*order_service.Order, error) {
	id := uuid.New()
	_, err := o.db.Exec(ctx, `INSERT INTO orders(id, external_id, type, customer_phone, customer_name, customer_id, status, to_address, to_location, discount_amount, amount, delivery_price, paid)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);`, id, req.ExternalId, req.Type, req.CustomerPhone, req.CustomerName, req.CustomerId, req.Status, req.ToAddress, req.ToLocation, req.DiscountAmount, req.Amount, req.DeliveryPrice, req.Paid)

	if err != nil {
		log.Println("failed to insert data to orders table: ", err)
		return nil, err
	}
	order, err := o.GetById(ctx, &order_service.OrderPrimaryKey{Id: id.String()})

	if err != nil {
		log.Println("failed to get data after insert data to orders table: ", err)
		return nil, err
	}
	return order, nil
}

func (o *orderRepo) GetById(ctx context.Context, req *order_service.OrderPrimaryKey) (*order_service.Order, error) {
	resp := &order_service.Order{}

	row := o.db.QueryRow(ctx, `SELECT * FROM orders WHERE id = $1;`, req.Id)
	err := row.Scan(&resp.Id, &resp.ExternalId, &resp.Type, &resp.CustomerPhone, &resp.CustomerName, &resp.CustomerId, &resp.Status, &resp.ToAddress, &resp.ToLocation, &resp.DiscountAmount, &resp.Amount, &resp.DeliveryPrice, &resp.Paid, &resp.CreatedAt, &resp.UpdatedAt, resp.DeletedAt)

	if err != nil {
		log.Println("failed to get a data from orders table: ", err)
		return nil, err
	}
	return resp, nil
}
