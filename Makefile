include .env

build:
	@go build -o bin/main.exe main.go

run:
	@go run main.go

live:
	@air

db-status:
	@goose -dir migrations $(GOOSE_DRIVER) $(GOOSE_DBSTRING) status

db-up:
	@goose -dir migrations $(GOOSE_DRIVER) $(GOOSE_DBSTRING) up

db-down:
	@goose -dir migrations $(GOOSE_DRIVER) $(GOOSE_DBSTRING) down