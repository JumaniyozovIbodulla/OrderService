package postgres

import (
	"context"
	"fmt"
	"order/genproto/order_service"
	"order/storage"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkt"
)

func toWKT(polygon *order_service.Polygon) (string, error) {
	n := len(polygon.Points)
	if n < 3 {
		return "", fmt.Errorf("polygon requires at least three points to form a valid shape")
	}

	points := make([]string, n+1)
	for i, point := range polygon.Points {
		points[i] = fmt.Sprintf("%f %f", point.Latitude, point.Longitude)
	}
	points[n] = fmt.Sprintf("%f %f", polygon.Points[0].Latitude, polygon.Points[0].Longitude)

	return fmt.Sprintf("POLYGON((%s))", strings.Join(points, ", ")), nil
}

func mapOrderTypeToPostgreSQL(orderType order_service.TypeEnum) string {
	switch orderType {
	case order_service.TypeEnum_self_pickup:
		return "self_pickup"
	case order_service.TypeEnum_delivery:
		return "delivery"
	default:
		return ""
	}
}

func mapPaymentEnumToPostgreSQL(paymentType order_service.PaymentEnum) string {
	switch paymentType {
	case order_service.PaymentEnum_waiting_for_payment:
		return "waiting_for_payment"
	case order_service.PaymentEnum_collecting:
		return "collecting"
	case order_service.PaymentEnum_shipping:
		return "shipping"
	case order_service.PaymentEnum_waiting_on_branch:
		return "waiting_on_branch"
	case order_service.PaymentEnum_finished:
		return "finished"
	case order_service.PaymentEnum_cancelled:
		return "cancelled"
	default:
		return ""
	}
}

func mapPostgreSQLToOrderType(orderType string) order_service.TypeEnum {
	switch orderType {
	case "self_pickup":
		return order_service.TypeEnum_self_pickup
	case "delivery":
		return order_service.TypeEnum_delivery
	default:
		return order_service.TypeEnum(0)
	}
}

func mapPostgreSQLToPaymentEnum(paymentStatus string) order_service.PaymentEnum {
	switch paymentStatus {
	case "waiting_for_payment":
		return order_service.PaymentEnum_waiting_for_payment
	case "collecting":
		return order_service.PaymentEnum_collecting
	case "shipping":
		return order_service.PaymentEnum_shipping
	case "waiting_on_branch":
		return order_service.PaymentEnum_waiting_on_branch
	case "finished":
		return order_service.PaymentEnum_finished
	case "cancelled":
		return order_service.PaymentEnum_cancelled
	default:
		return order_service.PaymentEnum(0)
	}
}

func fromWKT(wktString string) (*order_service.Polygon, error) {
	g, err := wkt.Unmarshal(wktString)
	if err != nil {
		return nil, err
	}

	switch g := g.(type) {
	case *geom.Polygon:
		var points []*order_service.Point
		for _, ring := range g.Coords() {
			for _, coord := range ring {
				point := &order_service.Point{
					Latitude:  coord.Y(),
					Longitude: coord.X(),
				}
				points = append(points, point)
			}
		}
		return &order_service.Polygon{Points: points}, nil
	default:
		return nil, fmt.Errorf("invalid WKT: expected Polygon")
	}
}

type orderRepo struct {
	db *pgxpool.Pool
}

func NewOrderRepo(db *pgxpool.Pool) storage.OrderRepo {
	return &orderRepo{
		db: db,
	}
}

func (o *orderRepo) Create(ctx context.Context, req *order_service.CreateOrder) (resp *order_service.Order, err error) {
	id := uuid.New()

	orderType := mapOrderTypeToPostgreSQL(req.Type)
	paymentStatus := mapPaymentEnumToPostgreSQL(req.Status)

	toLocationWKT, err := toWKT(req.ToLocation)

	if err != nil {
		return
	}
	_, err = o.db.Exec(ctx, `INSERT INTO orders(id, external_id, type, customer_phone, customer_name, customer_id, status, to_address, to_location, discount_amount, amount, delivery_price, paid)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, ST_GeomFromText($9, 4326), $10, $11, $12, $13);`, id, req.ExternalId, orderType, req.CustomerPhone, req.CustomerName, req.CustomerId, paymentStatus, req.ToAddress, toLocationWKT, req.DiscountAmount, req.Amount, req.DeliveryPrice, req.Paid)

	if err != nil {
		return
	}
	resp, err = o.GetById(ctx, &order_service.OrderPrimaryKey{Id: id.String()})

	if err != nil {
		return
	}
	return
}

func (o *orderRepo) GetById(ctx context.Context, req *order_service.OrderPrimaryKey) (resp *order_service.Order, err error) {
	resp = &order_service.Order{}

	row := o.db.QueryRow(ctx, `
		SELECT id, external_id, type, customer_phone, customer_name, customer_id, status, to_address, ST_AsText(to_location), discount_amount, amount, delivery_price, paid, TO_CHAR(created_at,'YYYY-MM-DD HH24:MI:SS TZH:TZM'), TO_CHAR(updated_at,'YYYY-MM-DD HH24:MI:SS TZH:TZM'), deleted_at
		FROM orders
		WHERE id = $1;
	`, req.Id)

	var toLocationWKT string
	var orderType, paymentStatus pgtype.Text

	err = row.Scan(
		&resp.Id, &resp.ExternalId, &orderType, &resp.CustomerPhone, &resp.CustomerName, &resp.CustomerId,
		&paymentStatus, &resp.ToAddress, &toLocationWKT, &resp.DiscountAmount, &resp.Amount, &resp.DeliveryPrice,
		&resp.Paid, &resp.CreatedAt, &resp.UpdatedAt, &resp.DeletedAt,
	)

	if err != nil {
		return
	}

	resp.Type = mapPostgreSQLToOrderType(orderType.String)
	resp.Status = mapPostgreSQLToPaymentEnum(paymentStatus.String)

	resp.ToLocation, err = fromWKT(toLocationWKT)
	if err != nil {
		return
	}
	return
}

func (o *orderRepo) Update(ctx context.Context, req *order_service.UpdateOrder) (resp *order_service.Order, err error) {

	orderType := mapOrderTypeToPostgreSQL(req.Type)
	paymentStatus := mapPaymentEnumToPostgreSQL(req.Status)

	toLocationWKT, err := toWKT(req.ToLocation)

	if err != nil {
		return
	}

	_, err = o.db.Exec(ctx, `
	UPDATE
		orders
	SET
		external_id = $2, type = $3, customer_phone = $4, customer_name = $5, customer_id = $6, status = $7, to_address = $8, to_location = $9, discount_amount = $10, amount = $11, delivery_price = $12, paid = $13, updated_at = NOW(), deleted_at = $14
	WHERE
		id = $1;`, req.Id, req.ExternalId, orderType, req.CustomerPhone, req.CustomerName, req.CustomerId, paymentStatus, req.ToAddress, toLocationWKT, req.DiscountAmount, req.Amount, req.DeliveryPrice, req.Paid, req.DeletedAt)

	if err != nil {
		return
	}

	resp, err = o.GetById(ctx, &order_service.OrderPrimaryKey{Id: req.Id})

	if err != nil {
		return
	}
	return
}

func (o *orderRepo) Delete(ctx context.Context, req *order_service.OrderPrimaryKey) (resp *order_service.Empty, err error) {
	_, err = o.db.Exec(ctx, `UPDATE orders SET deleted_at = EXTRACT(EPOCH FROM NOW()) WHERE id = $1;`, req.Id)

	if err != nil {
		return
	}
	return
}

func (o *orderRepo) GetAll(ctx context.Context, req *order_service.GetListOrderRequest) (resp *order_service.GetListOrderResponse, err error) {
	resp = &order_service.GetListOrderResponse{}
	filter := ""

	if req.Search != "" {
		filter = ` AND customer_name ILIKE '%` + req.Search + `%' `
	}

	rows, err := o.db.Query(ctx, ` 
	SELECT
	id, external_id, type, customer_phone, customer_name, customer_id, status, to_address, ST_AsText(to_location), discount_amount, amount, delivery_price, paid, TO_CHAR(created_at,'YYYY-MM-DD HH24:MI:SS TZH:TZM'), TO_CHAR(updated_at,'YYYY-MM-DD HH24:MI:SS TZH:TZM'), deleted_at
	FROM
		orders
	WHERE TRUE `+filter+` AND deleted_at = 0
	OFFSET
		$1
	LIMIT
		$2;`, req.Offset, req.Limit)

	if err != nil {
		return
	}

	for rows.Next() {
		var (
			order                    order_service.Order
			toLocationWKT            string
			orderType, paymentStatus pgtype.Text
		)

		if err = rows.Scan(
			&order.Id,
			&order.ExternalId,
			&orderType,
			&order.CustomerPhone,
			&order.CustomerName,
			&order.CustomerId,
			&paymentStatus,
			&order.ToAddress,
			&toLocationWKT,
			&order.DiscountAmount,
			&order.Amount,
			&order.DeliveryPrice,
			&order.Paid,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.DeletedAt); err != nil {
			return
		}

		order.Type = mapPostgreSQLToOrderType(orderType.String)
		order.Status = mapPostgreSQLToPaymentEnum(paymentStatus.String)

		order.ToLocation, err = fromWKT(toLocationWKT)
		if err != nil {
			return
		}

		resp.Orders = append(resp.Orders, &order)
	}

	err = o.db.QueryRow(ctx, `SELECT COUNT(*) FROM orders WHERE TRUE `+filter+` AND deleted_at = 0`).Scan(&resp.Count)
	if err != nil {
		return
	}
	return
}