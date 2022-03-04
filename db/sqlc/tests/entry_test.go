package tests

import (
	"context"
	"testing"
	"time"

	db "github.com/andrelsf/go-restapi/db/sqlc"
	"github.com/andrelsf/go-restapi/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, account db.Account) db.Entry {
	params := db.CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, params.AccountID, entry.AccountID)
	require.Equal(t, params.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntry(t, account)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	entryCreated := createRandomEntry(t, account)

	entryFound, err := testQueries.GetEntry(context.Background(), entryCreated.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entryFound)

	require.Equal(t, entryCreated.ID, entryFound.ID)
	require.Equal(t, entryCreated.AccountID, entryFound.AccountID)
	require.Equal(t, entryCreated.Amount, entryFound.Amount)
	require.WithinDuration(t, entryCreated.CreatedAt, entryFound.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	account := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntry(t, account)
	}

	params := db.ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	entries, err := testQueries.ListEntries(context.Background(), params)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, params.AccountID, entry.AccountID)
	}
}
