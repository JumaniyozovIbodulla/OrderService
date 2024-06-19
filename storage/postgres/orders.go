package postgres

import (
	"context"
	"fmt"
	"order/genproto/order_service"
	"order/storage"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
)

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

func mapPaymentTypeToPostgreSQL(paymentType order_service.PaymentType) string {
	switch paymentType {
	case order_service.PaymentType_uzum:
		return "uzum"
	case order_service.PaymentType_cash:
		return "cash"
	case order_service.PaymentType_terminal:
		return "terminal"
	default:
		return ""
	}
}

func mapPostgreSQLToPaymentType(paymentType string) order_service.PaymentType {
	switch paymentType {
	case "uzum":
		return order_service.PaymentType_uzum
	case "cash":
		return order_service.PaymentType_cash
	case "terminal":
		return order_service.PaymentType_terminal
	default:
		return order_service.PaymentType(0)
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
	
	var externalId string

	paymentType := mapPaymentTypeToPostgreSQL(req.PaymentType)
	orderType := mapOrderTypeToPostgreSQL(req.Type)
	paymentStatus := mapPaymentEnumToPostgreSQL(req.Status)

	row := o.db.QueryRow(ctx, `SELECT external_id FROM orders ORDER BY created_at DESC LIMIT 1;`)

	err := row.Scan(&externalId)
	if err != nil {
		_, err = o.db.Exec(ctx, `
			INSERT INTO 
				orders(id, external_id, type, customer_phone, customer_name, customer_id, payment_type, status, to_address, to_location, discount_amount, amount, delivery_price, paid, courier_id, courier_phone, courier_name)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, ST_SetSRID(ST_MakePoint($10, $11), 4326), $12, $13, $14, $15, $16, $17, $18);`, id, "num-000001", orderType, req.CourierPhone, req.CustomerName, req.CustomerId, paymentType, paymentStatus, req.ToAddress, req.ToLocation.Longitude, req.ToLocation.Latitude, req.DiscountAmount, req.Amount, req.DeliveryPrice, req.Paid, req.CourierId, req.CourierPhone, req.CourierName)
		if err != nil {
			return nil, err
		}
	} else {
		splitted := strings.Split(externalId, "-")
		number, err := strconv.Atoi(splitted[1])

		if err != nil {
			return nil, err
		}
		number++
		stringNumber := strconv.Itoa(number)
		length := 6 - len(stringNumber)
		zeros := strings.Repeat("0", length)
		result := fmt.Sprintf("num-%s%d", zeros, number)
		_, err = o.db.Exec(ctx, `
			INSERT INTO 
				orders(id, external_id, type, customer_phone, customer_name, customer_id, payment_type, status, to_address, to_location, discount_amount, amount, delivery_price, paid, courier_id, courier_phone, courier_name)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, ST_SetSRID(ST_MakePoint($10, $11), 4326), $12, $13, $14, $15, $16, $17, $18);`, id, result, orderType, req.CourierPhone, req.CustomerName, req.CustomerId, paymentType, paymentStatus, req.ToAddress, req.ToLocation.Longitude, req.ToLocation.Latitude, req.DiscountAmount, req.Amount, req.DeliveryPrice, req.Paid, req.CourierId, req.CourierPhone, req.CourierName)

		if err != nil {
			return nil, err
		}
	}

	resp, err := o.GetById(ctx, &order_service.OrderPrimaryKey{Id: id.String()})

	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *orderRepo) GetById(ctx context.Context, req *order_service.OrderPrimaryKey) (*order_service.Order, error) {
	resp := &order_service.Order{}

	row := o.db.QueryRow(ctx, `
		SELECT 
			id, 
			external_id, 
			type, 
			customer_phone, 
			customer_name, 
			customer_id, 
			payment_type,
			status, 
			to_address, 
			ST_Y(to_location) AS latitude, 
			ST_X(to_location) AS longitude, 
			discount_amount, 
			amount, 
			delivery_price, 
			paid, 
			courier_id,
			courier_phone,
			courier_name,
			TO_CHAR(created_at,'YYYY-MM-DD HH24:MI:SS TZH:TZM') AS created_at, 
			TO_CHAR(updated_at,'YYYY-MM-DD HH24:MI:SS TZH:TZM') AS updated_at
		FROM 
			orders
		WHERE 
			id = $1;
	`, req.Id)

	var (
		paymentType, orderType, paymentStatus pgtype.Text
		longitude, latitude                   float64
	)

	err := row.Scan(
		&resp.Id, &resp.ExternalId, &orderType, &resp.CustomerPhone, &resp.CustomerName, &resp.CustomerId, &paymentType,
		&paymentStatus, &resp.ToAddress, &latitude, &longitude, &resp.DiscountAmount, &resp.Amount, &resp.DeliveryPrice,
		&resp.Paid, &resp.CourierId, &resp.CourierPhone, &resp.CourierName, &resp.CreatedAt, &resp.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	resp.ToLocation = &order_service.Location{
		Latitude:  latitude,
		Longitude: longitude,
	}

	resp.PaymentType = mapPostgreSQLToPaymentType(paymentType.String)
	resp.Type = mapPostgreSQLToOrderType(orderType.String)
	resp.Status = mapPostgreSQLToPaymentEnum(paymentStatus.String)
	return resp, nil
}

func (o *orderRepo) Update(ctx context.Context, req *order_service.UpdateOrder) (*order_service.Order, error) {

	paymentType := mapPaymentTypeToPostgreSQL(req.PaymentType)
	orderType := mapOrderTypeToPostgreSQL(req.Type)
	paymentStatus := mapPaymentEnumToPostgreSQL(req.Status)

	_, err := o.db.Exec(ctx, `
	UPDATE
		orders
	SET
		external_id = $2, type = $3, customer_phone = $4, customer_name = $5, customer_id = $6, payment_type = $7, status = $8, to_address = $9, to_location = ST_SetSRID(ST_MakePoint($10, $11), 4326), discount_amount = $12, amount = $13, delivery_price = $14, paid = $15, courier_id = $16, courier_phone = $17, courier_name = $18, updated_at = NOW(), deleted_at = $19
	WHERE
		id = $1;`, req.Id, req.ExternalId, orderType, req.CustomerPhone, req.CustomerName, req.CustomerId, paymentType, paymentStatus, req.ToAddress, req.ToLocation.Longitude, req.ToLocation.Latitude, req.DiscountAmount, req.Amount, req.DeliveryPrice, req.Paid, req.CourierId, req.CourierPhone, req.CourierName, req.DeletedAt)

	if err != nil {
		return nil, err
	}

	resp, err := o.GetById(ctx, &order_service.OrderPrimaryKey{Id: req.Id})

	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *orderRepo) Delete(ctx context.Context, req *order_service.OrderPrimaryKey) (*order_service.Empty, error) {
	_, err := o.db.Exec(ctx, `UPDATE orders SET deleted_at = EXTRACT(EPOCH FROM NOW()) WHERE id = $1;`, req.Id)

	if err != nil {
		return nil, err
	}
	return &order_service.Empty{}, nil
}

func (o *orderRepo) GetAll(ctx context.Context, req *order_service.GetListOrderRequest) (*order_service.GetListOrderResponse, error) {
	resp := &order_service.GetListOrderResponse{}
	filter := ""

	if req.Search != "" {
		filter = ` AND customer_name ILIKE '%` + req.Search + `%' `
	}

	rows, err := o.db.Query(ctx, ` 
	SELECT 
		id, 
		external_id, 
		type, 
		customer_phone, 
		customer_name, 
		customer_id, 
		payment_type,
		status, 
		to_address, 
		ST_Y(to_location) AS latitude, 
		ST_X(to_location) AS longitude, 
		discount_amount, 
		amount, 
		delivery_price, 
		paid, 
		courier_id,
		courier_phone,
		courier_name,
		TO_CHAR(created_at,'YYYY-MM-DD HH24:MI:SS TZH:TZM') AS created_at, 
		TO_CHAR(updated_at,'YYYY-MM-DD HH24:MI:SS TZH:TZM') AS updated_at
	FROM
		orders
	WHERE TRUE `+filter+` AND deleted_at = 0
	OFFSET
		$1
	LIMIT
		$2;`, req.Offset, req.Limit)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			order                                 order_service.Order
			paymentType, orderType, paymentStatus pgtype.Text
			longitude, latitude                   float64
		)

		if err = rows.Scan(
			&order.Id,
			&order.ExternalId,
			&orderType,
			&order.CustomerPhone,
			&order.CustomerName,
			&order.CustomerId,
			&paymentType,
			&paymentStatus,
			&order.ToAddress,
			&latitude,
			&longitude,
			&order.DiscountAmount,
			&order.Amount,
			&order.DeliveryPrice,
			&order.Paid,
			&order.CourierId,
			&order.CourierPhone,
			&order.CourierName,
			&order.CreatedAt,
			&order.UpdatedAt); err != nil {
			return nil, err
		}

		order.PaymentType = mapPostgreSQLToPaymentType(paymentType.String)
		order.Type = mapPostgreSQLToOrderType(orderType.String)
		order.Status = mapPostgreSQLToPaymentEnum(paymentStatus.String)

		order.ToLocation = &order_service.Location{
			Latitude:  latitude,
			Longitude: longitude,
		}

		resp.Orders = append(resp.Orders, &order)
	}
	err = o.db.QueryRow(ctx, `SELECT COUNT(*) FROM orders WHERE TRUE `+filter+` AND deleted_at = 0`).Scan(&resp.Count)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
