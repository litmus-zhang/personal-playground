package bank

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/litmus-zhang/simple_bank/util"
	"github.com/stretchr/testify/require"
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

func CreateAccountForTest(t *testing.T) Account {
	user := CreateUserForTest(t)
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(1000),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}

func CreateUserForTest(t *testing.T) CreateUserRow {
	hash, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		FullName:       util.RandomOwner() + " " + util.RandomOwner(),
		Email:          util.RandomOwner() + "@test.com",
		HashedPassword: hash,
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.True(t, user.PasswordChangedAt.IsZero())
	return user
}
