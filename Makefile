build:
	go build -o bin/allocator cmd/main.go

run:
	go run cmd/main.go

up:
	goose -dir=migrations/ postgres "postgres://username:password@localhost:5432/alloc?sslmode=disable" up
