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
	dbSource = "postgresql://root:secret@localhost:5001/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testQueries = New(testDB)

	code := m.Run()

	testDB.Close()

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
