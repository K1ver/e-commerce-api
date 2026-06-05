.PHONY: run migrate-up migrate-down swagger build test

run:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

test:
	go test ./...

migrate-up:
	goose -dir migration postgres "host=localhost port=5432 user=postgres password=postgres dbname=ecommerce sslmode=disable" up

migrate-down:
	goose -dir migration postgres "host=localhost port=5432 user=postgres password=postgres dbname=ecommerce sslmode=disable" down

swagger:
	swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal
