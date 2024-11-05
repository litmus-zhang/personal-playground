package bank

import (
	"context"
	"testing"

	"github.com/litmus-zhang/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func CreateEntryForTest(t *testing.T) Entry {
	account := CreateAccountForTest(t)

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(10),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.ID)

	return entry
}
func TestCreateEntry(t *testing.T) {
	CreateEntryForTest(t)

}

func TestGetEntry(t *testing.T) {
	entry := CreateEntryForTest(t)

	entry2, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry.ID, entry2.ID)
	require.Equal(t, entry.AccountID, entry2.AccountID)
	require.Equal(t, entry.Amount, entry.Amount)
}

func TestGetAllEntries(t *testing.T) {
	for i := 0; i < 5; i++ {
		CreateEntryForTest(t)
	}
	arg := ListEntriesParams{
		Limit:  5,
		Offset: 0,
	}
	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)
}
