# RAGnarok

## Quick start

1. Set variables in `.env` file

2. Run commands:
```
docker-compose up -d
make seed
make gateway
```

3. Profit

## Migrations
for migrations use:
go install github.com/pressly/goose/v3/cmd/goose@latest

new migration:
goose create add_lorem_ipsum_table sql

## DB admin
http://127.0.0.1:5050
