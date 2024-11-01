PROTO_SRC_DIR := internal/grpc/proto
PROTO_OUT_DIR := pkg/pb
DB_URL := "postgres://username:password@localhost:5432/alloc?sslmode=disable"

.PHONY: proto-gen migrate-create migrate-up build run up

proto-gen:
	protoc --go_out=$(PROTO_OUT_DIR) \
	       --go_opt=paths=source_relative \
	       --go-grpc_out=$(PROTO_OUT_DIR) \
	       --go-grpc_opt=paths=source_relative \
	       --proto_path=$(PROTO_SRC_DIR) \
	       $(PROTO_SRC_DIR)/allocator_svc.proto

migrate-create:
	goose -dir=migrations create $(name) sql

migrate-up:
	goose -dir=migrations postgres $(DB_URL) up

build:
	go build -o bin/allocator cmd/main.go

run:
	go run cmd/main.go

up:
	goose -dir=migrations postgres $(DB_URL) up
