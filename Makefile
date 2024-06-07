CURRENT_DIR := $(shell pwd)

gen-proto:
	sudo rm -rf ${CURRENT_DIR}/genproto/order_service
	mkdir -p ${CURRENT_DIR}/genproto/order_service
	sudo rm -rf ${CURRENT_DIR}/genproto/order_product_service
	mkdir -p ${CURRENT_DIR}/genproto/order_product_service
	sudo rm -rf ${CURRENT_DIR}/genproto/order_notes
	mkdir -p ${CURRENT_DIR}/genproto/order_notes
	protoc --proto_path=proto/order_protos --go_out=${CURRENT_DIR}/genproto/order_service --go_opt=paths=source_relative --go-grpc_out=${CURRENT_DIR}/genproto/order_service --go-grpc_opt=paths=source_relative proto/order_protos/orders.proto
	protoc --proto_path=proto/order_protos --go_out=${CURRENT_DIR}/genproto/order_product_service --go_opt=paths=source_relative --go-grpc_out=${CURRENT_DIR}/genproto/order_product_service --go-grpc_opt=paths=source_relative proto/order_protos/order_products.proto
	protoc --proto_path=proto/order_protos --go_out=${CURRENT_DIR}/genproto/order_notes --go_opt=paths=source_relative --go-grpc_out=${CURRENT_DIR}/genproto/order_notes --go-grpc_opt=paths=source_relative proto/order_protos/order_status_notes.proto

swag_init:
	swag init -g api/main.go -o api/docs

run:
	go run cmd/main.go

git-push:
	git push origin main