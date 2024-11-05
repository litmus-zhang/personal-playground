package bank

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateTransferForTest(t *testing.T) Transfer {
	account1 := CreateAccountForTest(t)
	account2 := CreateAccountForTest(t)

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        10,
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	return transfer
}
func TestCreateTransfer(t *testing.T) {
	CreateAccountForTest(t)
}

func TestGetTransfer(t *testing.T) {
	transfer := CreateTransferForTest(t)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)
	require.Equal(t, transfer.ID, transfer2.ID)

}

func TestGetAllTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateTransferForTest(t)
	}
	arg := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}
	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)
}
