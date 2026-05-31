.PHONY: migrate-up migrate-down sqlc run-api run-worker build tidy

DATABASE_URL ?= postgres://postgres:Admin123@localhost:5432/wa_dashboard?sslmode=disable
MIGRATIONS_PATH ?= db/migrations

migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down 1

sqlc:
	sqlc generate

run-api:
	go run ./cmd/api

run-worker:
	go run ./cmd/worker

build:
	go build ./...

tidy:
	go mod tidy
