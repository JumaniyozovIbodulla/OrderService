package service

import (
	"context"
	on "order/genproto/order_notes"

	"github.com/saidamir98/udevs_pkg/logger"
)

func (o *OrderNotesService) Create(ctx context.Context, req *on.CreateOrderNotes) (resp *on.OrderNotes, err error) {
	o.log.Info("Create Note: ", logger.Any("req", req))

	resp, err = o.strg.OrderNotes().Create(ctx, req)

	if err != nil {
		o.log.Error("Create Note: ", logger.Error(err))
		return
	}
	return
}

func (o *OrderNotesService) GetById(ctx context.Context, req *on.OrderNotesPrimaryKey) (resp *on.OrderNotes, err error) {
	o.log.Info("Get Note: ", logger.Any("req", req))

	resp, err = o.strg.OrderNotes().GetById(ctx, req)

	if err != nil {
		o.log.Error("Get Note: ", logger.Error(err))
		return
	}
	return
}


func (o *OrderNotesService) Update(ctx context.Context, req *on.UpdateOrderNotes) (resp *on.OrderNotes, err error) {
	o.log.Info("Update Note: ", logger.Any("req", req))

	resp, err = o.strg.OrderNotes().Update(ctx, req)

	if err != nil {
		o.log.Error("Update Note: ", logger.Error(err))
		return
	}
	return
}


func (o *OrderNotesService) Delete(ctx context.Context, req *on.OrderNotesPrimaryKey) (resp *on.Empty, err error) {
	o.log.Info("Delete Note: ", logger.Any("req", req))

	resp, err = o.strg.OrderNotes().Delete(ctx, req)

	if err != nil {
		o.log.Error("Delete Note: ", logger.Error(err))
		return
	}
	return
}


func (o *OrderNotesService) GetAll(ctx context.Context, req *on.GetListOrderNotesRequest) (resp *on.GetListOrderNotesResponse, err error) {
	o.log.Info("Get all notes: ", logger.Any("req", req))

	resp, err = o.strg.OrderNotes().GetAll(ctx, req)

	if err != nil {
		o.log.Error("Get all notes: ", logger.Error(err))
		return
	}
	return
}