startdb:
	docker compose up -d

stopdb:
	docker compose down

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5001/simple_bank?sslmode=disable"  -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5001/simple_bank?sslmode=disable"  -verbose up 1

migrateup_test:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5001/postgres?sslmode=disable"  -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5001/simple_bank?sslmode=disable"  -verbose down
migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5001/simple_bank?sslmode=disable"  -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

run_tests:
	make migrateup && make test

mock:
	mockgen -package mockdb -destination db/mock/store.go  github.com/litmus-zhang/simple_bank/bank Store

.PHONY: startdb stopdb migrateup migrateup1 migrateup_test migratedown migratedown1 sqlc  test run_tests server mock