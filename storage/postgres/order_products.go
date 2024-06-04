package postgres

import (
	"context"
	orp "order/genproto/order_product_service"
	"order/storage"

	"github.com/google/uuid"
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
	resp = &orp.OrderProduct{}
	id := uuid.New()
	_, err = o.db.Exec(ctx, `
	INSERT INTO
		order_products(id, product_id, count, discount_price, price, order_id)
	VALUES($1, $2, $3, $4, $5, $6);`, id, req.ProductId, req.Count, req.DiscountPrice, req.Price, req.OrderId)

	if err != nil {
		return
	}
	resp, err = o.GetById(ctx, &orp.OrderProductPrimaryKey{Id: id.String()})
	return
}

func (o *orderProducsRepo) GetById(ctx context.Context, req *orp.OrderProductPrimaryKey) (resp *orp.OrderProduct, err error) {
	resp = &orp.OrderProduct{}

	row := o.db.QueryRow(ctx, `SELECT id, product_id, count, discount_price, price, order_id, TO_CHAR(created_at,'YYYY-MM-DD HH24:MI:SS TZH:TZM'), TO_CHAR(updated_at,'YYYY-MM-DD HH24:MI:SS TZH:TZM'), deleted_at FROM order_products WHERE id = $1;`, req.Id)
	err = row.Scan(&resp.Id, &resp.ProductId, &resp.Count, &resp.DiscountPrice, &resp.Price, &resp.OrderId, &resp.CreatedAt, &resp.UpdatedAt, &resp.DeletedAt)
	return
}

func (o *orderProducsRepo) GetAll(ctx context.Context, req *orp.GetListOrderProductRequest) (resp *orp.GetListOrderProductResponse, err error) {
	resp = &orp.GetListOrderProductResponse{}

	rows, err := o.db.Query(ctx, `SELECT id, product_id, count, discount_price, price, order_id, TO_CHAR(created_at,'YYYY-MM-DD HH24:MI:SS TZH:TZM'), TO_CHAR(updated_at,'YYYY-MM-DD HH24:MI:SS TZH:TZM'), deleted_at FROM order_products WHERE deleted_at = 0  OFFSET $1 LIMIT $2;`, req.Offset, req.Limit)

	if err != nil {
		return
	}

	for rows.Next() {
		var (
			order_product orp.OrderProduct
		)

		if err = rows.Scan(&order_product.Id, &order_product.ProductId, &order_product.Count, &order_product.DiscountPrice, &order_product.Price, &order_product.OrderId, &order_product.CreatedAt, &order_product.UpdatedAt, &order_product.DeletedAt); err != nil {
			return
		}

		resp.OrderProducts = append(resp.OrderProducts, &order_product)
	}

	err = o.db.QueryRow(ctx, `SELECT COUNT(*) FROM order_products WHERE deleted_at = 0`).Scan(&resp.Count)
	if err != nil {
		return
	}
	return
}

func (o *orderProducsRepo) Update(ctx context.Context, req *orp.UpdateOrderProduct) (resp *orp.OrderProduct, err error) {
	resp = &orp.OrderProduct{}
	_, err = o.db.Exec(ctx, `
	UPDATE 
		order_products
	SET
		product_id = $2, count = $3, discount_price = $4, price = $5, order_id = $6, updated_at = NOW(), deleted_at = $7
	WHERE
		id = $1;`, req.Id, req.ProductId, req.Count, req.DiscountPrice, req.Price, req.OrderId, req.DeletedAt)

	if err != nil {
		return
	}
	resp, err = o.GetById(ctx, &orp.OrderProductPrimaryKey{Id: req.Id})
	return
}

func (o *orderProducsRepo) Delete(ctx context.Context, req *orp.OrderProductPrimaryKey) (resp *orp.Empty, err error) {
	_, err = o.db.Exec(ctx, `UPDATE order_products SET deleted_at = EXTRACT(EPOCH FROM NOW()) WHERE id = $1;`, req.Id)
	return
}
