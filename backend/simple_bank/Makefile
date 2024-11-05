startdb:
	docker compose up -d

stopdb:
	docker compose down

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5001/simple_bank?sslmode=disable"  -verbose up

migrateup_test:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5001/postgres?sslmode=disable"  -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5001/simple_bank?sslmode=disable"  -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

run_test:
	make migrateup_test && make  test

.PHONY: startdb stopdb migrateup migrateup_test migratedown sqlc  test run_test