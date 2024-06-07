package postgres

import (
	"context"
	orn "order/genproto/order_notes"
	"order/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
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

func PaymentEnumToPostgreSQL(paymentType orn.PaymentEnum) string {
	switch paymentType {
	case orn.PaymentEnum_waiting_for_payment:
		return "waiting_for_payment"
	case orn.PaymentEnum_collecting:
		return "collecting"
	case orn.PaymentEnum_shipping:
		return "shipping"
	case orn.PaymentEnum_waiting_on_branch:
		return "waiting_on_branch"
	case orn.PaymentEnum_finished:
		return "finished"
	case orn.PaymentEnum_cancelled:
		return "cancelled"
	default:
		return ""
	}
}

func PostgreSQLToPaymentEnum(paymentStatus string) orn.PaymentEnum {
	switch paymentStatus {
	case "waiting_for_payment":
		return orn.PaymentEnum_waiting_for_payment
	case "collecting":
		return orn.PaymentEnum_collecting
	case "shipping":
		return orn.PaymentEnum_shipping
	case "waiting_on_branch":
		return orn.PaymentEnum_waiting_on_branch
	case "finished":
		return orn.PaymentEnum_finished
	case "cancelled":
		return orn.PaymentEnum_cancelled
	default:
		return orn.PaymentEnum(0)
	}
}

func (o *orderNotesRepo) Create(ctx context.Context, req *orn.CreateOrderNotes) (*orn.OrderNotes, error) {
	id := uuid.New()
	paymentStatus := PaymentEnumToPostgreSQL(req.Status)
	_, err := o.db.Exec(ctx, `
	INSERT INTO
		order_status_notes(id, order_id, status, user_id, reason)	
	VALUES($1, $2, $3, $4, $5);`, id, req.OrderId, paymentStatus, req.UserId, req.Reason)

	if err != nil {
		return nil, err
	}
	resp, err := o.GetById(ctx, &orn.OrderNotesPrimaryKey{Id: id.String()})
	
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *orderNotesRepo) GetById(ctx context.Context, req *orn.OrderNotesPrimaryKey) (*orn.OrderNotes, error) {
	resp := &orn.OrderNotes{}
	var	paymentStatus pgtype.Text

	row := o.db.QueryRow(ctx, `SELECT id, order_id, status, user_id, reason, TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI:SS TZH:TZM') FROM order_status_notes WHERE id = $1;`, req.Id)
	err := row.Scan(&resp.Id, &resp.OrderId, &paymentStatus, &resp.UserId, &resp.Reason, &resp.CreatedAt)
	if  err != nil {
		return nil, err
	}
	resp.Status = PostgreSQLToPaymentEnum(paymentStatus.String)
	return resp, nil
}

func (o *orderNotesRepo) GetAll(ctx context.Context, req *orn.GetListOrderNotesRequest) (*orn.GetListOrderNotesResponse, error) {
	resp := &orn.GetListOrderNotesResponse{}
	rows, err := o.db.Query(ctx, `SELECT id, order_id, status, user_id, reason, TO_CHAR(created_at,'YYYY-MM-DD HH24:MI:SS TZH:TZM') FROM order_status_notes OFFSET $1 LIMIT $2;`, req.Offset, req.Limit)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			order_notes   orn.OrderNotes
			paymentStatus pgtype.Text
		)

		err = rows.Scan(&order_notes.Id, &order_notes.OrderId, &paymentStatus, &order_notes.UserId, &order_notes.Reason, &order_notes.CreatedAt)
		if err != nil {
			return nil, err
		}

		order_notes.Status = PostgreSQLToPaymentEnum(paymentStatus.String)
		resp.OrderNotes = append(resp.OrderNotes, &order_notes)
	}

	err = o.db.QueryRow(ctx, `SELECT COUNT(*) FROM order_status_notes;`).Scan(&resp.Count)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (o *orderNotesRepo) Update(ctx context.Context, req *orn.UpdateOrderNotes) (*orn.OrderNotes, error) {
	resp := &orn.OrderNotes{}
	paymentStatus := PaymentEnumToPostgreSQL(req.Status)
	_, err := o.db.Exec(ctx, `UPDATE order_status_notes SET order_id = $2, status = $3, user_id = $4, reason = $5 WHERE id = $1;`, req.Id, req.OrderId, paymentStatus, req.UserId, req.Reason)

	if err != nil {
		return nil, err
	}
	resp, err = o.GetById(ctx, &orn.OrderNotesPrimaryKey{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *orderNotesRepo) Delete(ctx context.Context, req *orn.OrderNotesPrimaryKey) (*orn.Empty, error) {
	resp := &orn.Empty{}
	_, err := o.db.Exec(ctx, `DELETE FROM order_status_notes WHERE id = $1;`, req.Id)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
