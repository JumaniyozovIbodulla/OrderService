package client

import "order/config"

type OrderServiceManager interface{}

type OrderProductsManager interface{}

type OrderNotesManager interface{}

type grpcOrderClients struct{}

type grpcOrderProductClients struct{}

type grpcOrderNotesClients struct{}

func NewGrpcOrderClients(cfg config.Config) (OrderServiceManager, error) {
	return *&grpcOrderClients{}, nil
}

func NewGrpcOrderProductClients(cfg config.Config) (OrderProductsManager, error) {
	return *&grpcOrderProductClients{}, nil
}

func NewGrpcOrderNotesClients(cfg config.Config) (OrderNotesManager, error) {
	return *&grpcOrderNotesClients{}, nil
}