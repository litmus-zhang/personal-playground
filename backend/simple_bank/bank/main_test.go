package bank

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5002/simple_bank_test?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testQueries = New(conn)

	code := m.Run()

	conn.Close()

	os.Exit(code)
}

func cleanDB(t *testing.T) {
	_, err := testQueries.db.ExecContext(context.Background(), `
	TRUNCATE transfer CASCADE;
	TRUNCATE entries CASCADE;
	TRUNCATE accounts CASCADE;
	`)
	if err != nil {
		t.Fatal("cannot truncate db:", err)
	}

}
