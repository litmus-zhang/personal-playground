package bank

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/litmus-zhang/simple_bank/util"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	config, err := util.LoadConfig("./..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBURL)
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
