package postgres

import (
	"context"
	"order/genproto/order_notes"
	"order/storage"

	"github.com/jackc/pgx/v4/pgxpool"
)

type orderNotesRepo struct {
	db *pgxpool.Pool
}

func NewOrderNotesRepo(db *pgxpool.Pool) storage.OrderNotesRepo {
	return &orderNotesRepo{
		db: db,
	}
}


func (o *orderNotesRepo) Create(ctx context.Context, req *order_notes.CreateOrderNotes) (resp *order_notes.OrderNotes, err error) {
	return
}