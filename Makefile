include .env

cli:
	@go build -C cli -o ../bin/cli .

scraper_service:
	@go build -C scraper_service -o ../bin/scraper_service .
	@./bin/scraper_service

indexer_service:
	@go build -C indexer_service -o ../bin/indexer_service .
	@./bin/indexer_service

gateway:
	@go build -C gateway -o ../bin/gateway .
	@./bin/gateway

db-status:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING="postgresql://$(DB_USER):$(DB_PASSWORD)@127.0.0.1:5432/$(DB_NAME)" GOOSE_MIGRATION_DIR="db/migrations" goose status

db-up:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING="postgresql://$(DB_USER):$(DB_PASSWORD)@127.0.0.1:5432/$(DB_NAME)" GOOSE_MIGRATION_DIR="db/migrations" goose up

db-down:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING="postgresql://$(DB_USER):$(DB_PASSWORD)@127.0.0.1:5432/$(DB_NAME)" GOOSE_MIGRATION_DIR="db/migrations" goose down

models:
	@PSQL_HOST=$(DB_HOST) PSQL_DBNAME=$(DB_NAME) PSQL_USER=$(DB_USER) PSQL_PASS=$(DB_PASSWORD) PSQL_SSLMODE=disable sqlboiler psql --wipe

seed:
	@go run scripts/seed.go

.PHONY: cli scraper_service indexer_service gateway models
