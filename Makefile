CURRENT_DIR := $(shell pwd)

gen-proto:
	sudo rm -rf ${CURRENT_DIR}/genproto/order_service
	mkdir -p ${CURRENT_DIR}/genproto/order_service
	protoc --proto_path=protos/order_service --go_out=${CURRENT_DIR}/genproto/order_service --go_opt=paths=source_relative --go-grpc_out=${CURRENT_DIR}/genproto/order_service --go-grpc_opt=paths=source_relative protos/order_service/orders.proto

swag_init:
	swag init -g api/main.go -o api/docs

run:
	go run cmd/main.go
	
git-push:
	git push origin main