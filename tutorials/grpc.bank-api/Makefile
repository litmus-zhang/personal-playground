project_name = "micro_bank"
DB_URL=postgresql://root:root@localhost:5432/$(project_name)?sslmode=disable

startdb:
	@echo "Starting Postgres..."
	brew services start postgresql # for mac
	docker compose up -d db # for other platform with docker installed
	psql postgres
	CREATE USER root WITH SUPERUSER PASSWORD 'root';
	CREATE DATABASE micro_bank;

stopdb:
	docker compose down -v --remove-orphans --rmi all

migrateup:
	migrate -path internal/db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path internal/db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path internal/db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path internal/db/migration -database "$(DB_URL)" -verbose down 1

new_migration:
	migrate create -ext sql -dir internal/db/migration -seq $(name)

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

server:
	go run cmd/main.go

mock:
	mockgen -package mockdb -destination internal/db/mock/store.go github.com/litmus-zhang/task_manager/internal/db Store

.PHONY: startdb stopdb migrateup migrateup1 migratedown migratedown1 new_migration sqlc test server mock