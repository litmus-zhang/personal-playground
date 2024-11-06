package bank

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := CreateAccountForTest(t)
	account2 := CreateAccountForTest(t)

	fmt.Print("Before transfer\n")
	fmt.Printf("Account 1: %v\n", account1.Balance)
	fmt.Printf("Account 2: %v\n", account2.Balance)

	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	n := 5
	for i := 0; i < n; i++ {

		go func() {

			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result

		}()
	}
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer

		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)

		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		// Check the balance of the accounts after the transfer
		fromAccount := result.FromAccount

		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)
		// require.Equal(t, account1.Balance-amount, fromAccount.Balance)

		toAccount := result.ToAccount

		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)
		// require.Equal(t, account2.Balance+amount, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

	}

	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount1)
	require.Equal(t, account1.ID, updatedAccount1.ID)
	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount2)
	require.Equal(t, account2.ID, updatedAccount2.ID)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)

	fmt.Print("After transfer\n")
	fmt.Printf("Account 1: %v\n", account1.Balance)
	fmt.Printf("Account 2: %v\n", account2.Balance)
}

func TestTransferTxForDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := CreateAccountForTest(t)
	account2 := CreateAccountForTest(t)

	fmt.Print("Before transfer\n")
	fmt.Printf("Account 1: %v\n", account1.Balance)
	fmt.Printf("Account 2: %v\n", account2.Balance)

	amount := int64(10)

	errs := make(chan error)

	n := 10
	for i := 0; i < n; i++ {
		FromAccountID := account1.ID
		ToAccountID := account2.ID

		if i%2 == 0 {
			FromAccountID = account2.ID
			ToAccountID = account1.ID
		}

		go func() {

			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: FromAccountID,
				ToAccountID:   ToAccountID,
				Amount:        amount,
			})
			errs <- err

		}()
	}
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

	}

	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.Equal(t, account1.Balance, updatedAccount1.Balance)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)

	fmt.Print("After transfer\n")
	fmt.Printf("Account 1: %v\n", account1.Balance)
	fmt.Printf("Account 2: %v\n", account2.Balance)
}
