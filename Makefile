# Makefile
migrate-down:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" goose -dir ./migrations down

migrate-up:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" goose -dir ./migrations up

gen-docs:
	swag init -g cmd/main.go -o docs
