package service

import (
	"context"
	opc "order/genproto/order_product_service"

	"github.com/saidamir98/udevs_pkg/logger"
)

func (o *OrderProductService) Create(ctx context.Context, req *opc.CreateOrderProduct) (*opc.OrderProduct, error) {
	o.log.Info("Create OrderProduct: ", logger.Any("req", req))

	resp, err := o.strg.OrderProduct().Create(ctx, req)

	if err != nil {
		o.log.Error("Create OrderProduct: ", logger.Error(err))
		return nil, err
	}
	return resp, nil
}

func (o *OrderProductService) GetById(ctx context.Context, req *opc.OrderProductPrimaryKey) (*opc.OrderProduct, error) {
	o.log.Info("Get OrderProduct: ", logger.Any("req", req))

	resp, err := o.strg.OrderProduct().GetById(ctx, req)

	if err != nil {
		o.log.Error("Get OrderProduct: ", logger.Error(err))
		return nil, err
	}
	return resp, nil
}


func (o *OrderProductService) Update(ctx context.Context, req *opc.UpdateOrderProduct) (*opc.OrderProduct, error) {
	o.log.Info("Update OrderProduct: ", logger.Any("req", req))

	resp, err := o.strg.OrderProduct().Update(ctx, req)

	if err != nil {
		o.log.Error("Update OrderProduct: ", logger.Error(err))
		return nil, err
	}
	return resp, nil
}


func (o *OrderProductService) Delete(ctx context.Context, req *opc.OrderProductPrimaryKey) (*opc.Empty, error) {
	o.log.Info("Delete OrderProduct: ", logger.Any("req", req))

	resp, err := o.strg.OrderProduct().Delete(ctx, req)

	if err != nil {
		o.log.Error("Delete OrderProduct: ", logger.Error(err))
		return nil, err
	}
	return resp, nil
}

func (o *OrderProductService) GetAll(ctx context.Context, req *opc.GetListOrderProductRequest) (*opc.GetListOrderProductResponse, error) {
	o.log.Info("GetAll OrderProduct: ", logger.Any("req", req))

	resp, err := o.strg.OrderProduct().GetAll(ctx, req)

	if err != nil {
		o.log.Error("GetAll OrderProduct: ", logger.Error(err))
		return nil, err
	}
	return resp, nil
}