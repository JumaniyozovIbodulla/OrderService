package client

import "order/config"

type OrderServiceManager interface{}

type grpcClients struct{}

func NewGrpcClients(cfg config.Config) (OrderServiceManager, error) {
	return *&grpcClients{}, nil
}