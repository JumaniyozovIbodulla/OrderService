package service

import (
	"context"
	on "order/genproto/order_notes"

	"github.com/saidamir98/udevs_pkg/logger"
)

func (o *OrderNotesService) Create(ctx context.Context, req *on.CreateOrderNotes) (*on.OrderNotes, error) {
	o.log.Info("Create Note: ", logger.Any("req", req))

	resp, err := o.strg.OrderNotes().Create(ctx, req)

	if err != nil {
		o.log.Error("Create Note: ", logger.Error(err))
		return nil, err
	}
	return resp, nil
}

func (o *OrderNotesService) GetById(ctx context.Context, req *on.OrderNotesPrimaryKey) (*on.OrderNotes, error) {
	o.log.Info("Get Note: ", logger.Any("req", req))

	resp, err := o.strg.OrderNotes().GetById(ctx, req)

	if err != nil {
		o.log.Error("Get Note: ", logger.Error(err))
		return nil, err
	}
	return resp, nil
}


func (o *OrderNotesService) Update(ctx context.Context, req *on.UpdateOrderNotes) (*on.OrderNotes, error) {
	o.log.Info("Update Note: ", logger.Any("req", req))

	resp, err := o.strg.OrderNotes().Update(ctx, req)

	if err != nil {
		o.log.Error("Update Note: ", logger.Error(err))
		return nil, err
	}
	return resp, nil
}


func (o *OrderNotesService) Delete(ctx context.Context, req *on.OrderNotesPrimaryKey) (*on.Empty, error) {
	o.log.Info("Delete Note: ", logger.Any("req", req))

	resp, err := o.strg.OrderNotes().Delete(ctx, req)

	if err != nil {
		o.log.Error("Delete Note: ", logger.Error(err))
		return nil, err
	}
	return resp, nil
}


func (o *OrderNotesService) GetAll(ctx context.Context, req *on.GetListOrderNotesRequest) (*on.GetListOrderNotesResponse, error) {
	o.log.Info("Get all notes: ", logger.Any("req", req))

	resp, err := o.strg.OrderNotes().GetAll(ctx, req)

	if err != nil {
		o.log.Error("Get all notes: ", logger.Error(err))
		return nil, err
	}
	return resp, nil
}