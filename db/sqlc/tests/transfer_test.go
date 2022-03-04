package tests

import (
	"context"
	"testing"
	"time"

	db "github.com/andrelsf/go-restapi/db/sqlc"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, fromAccount, toAccount db.Account) db.Transfer {
	params := db.CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, params.FromAccountID, transfer.FromAccountID)
	require.Equal(t, params.ToAccountID, transfer.ToAccountID)
	require.Equal(t, params.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	createRandomTransfer(t, fromAccount, toAccount)
}

func TestGetTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	transferRegistered := createRandomTransfer(t, fromAccount, toAccount)
	transferFound, err := testQueries.GetTransfer(context.Background(), transferRegistered.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transferFound)

	require.Equal(t, transferRegistered.ID, transferFound.ID)
	require.Equal(t, transferRegistered.FromAccountID, transferFound.FromAccountID)
	require.Equal(t, transferRegistered.ToAccountID, transferFound.ToAccountID)
	require.Equal(t, transferRegistered.Amount, transferFound.Amount)
	require.WithinDuration(t, transferRegistered.CreatedAt, transferFound.CreatedAt, time.Second)
}

func TestListTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	for i := 0; i < 5; i++ {
		createRandomTransfer(t, fromAccount, toAccount)
		createRandomTransfer(t, toAccount, fromAccount)
	}

	params := db.ListTransfersParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   fromAccount.ID,
		Limit:         5,
		Offset:        5,
	}

	tranfers, err := testQueries.ListTransfers(context.Background(), params)
	require.NoError(t, err)
	require.Len(t, tranfers, 5)

	for _, transfer := range tranfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == fromAccount.ID || transfer.ToAccountID == fromAccount.ID)
	}
}
