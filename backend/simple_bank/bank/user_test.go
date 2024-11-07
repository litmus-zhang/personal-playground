package bank

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	CreateUserForTest(t)
}

func TestGetUser(t *testing.T) {
	user := CreateUserForTest(t)
	u, err := testQueries.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, u)
	require.Equalf(t, u.Username, user.Username, "account id should be the same")
	require.Equal(t, u.FullName, user.FullName)
	require.Equal(t, u.Email, user.Email)
	require.WithinDuration(t, u.CreatedAt, user.CreatedAt, time.Second)
}
