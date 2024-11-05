startdb:
	docker compose up -d
stopdb:
	docker compose down

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5001/simple_bank?sslmode=disable"  -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5001/simple_bank?sslmode=disable"  -verbose down

sqlc:
	sqlc generate

.PHONY: startdb stopdb migrateup migratedown sqlc