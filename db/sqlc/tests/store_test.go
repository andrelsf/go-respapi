package tests

import (
	"context"
	"fmt"
	"testing"

	db "github.com/andrelsf/go-restapi/db/sqlc"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := db.NewStore(testDB)

	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	fmt.Println(">> before:", fromAccount.Balance, toAccount.Balance)

	// run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	// Channels
	errs := make(chan error)
	results := make(chan db.TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), db.TransferTxParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	// Check results
	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// Check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, fromAccount.ID, transfer.FromAccountID)
		require.Equal(t, toAccount.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// Check Entries
		// fromEntry
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)

		require.Equal(t, fromAccount.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// toEntry
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)

		require.Equal(t, toAccount.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// Check Accounts
		resultFromAccount := result.FromAccount
		require.NotEmpty(t, resultFromAccount)
		require.Equal(t, fromAccount.ID, resultFromAccount.ID)

		resultToAccount := result.ToAccount
		require.NotEmpty(t, resultToAccount)
		require.Equal(t, toAccount.ID, resultToAccount.ID)

		// check accounts' balance
		fmt.Println(">> tx:", resultFromAccount.Balance, resultToAccount.Balance)
		diffBalanceFromAccount := fromAccount.Balance - resultFromAccount.Balance
		diffBalanceToAccount := resultToAccount.Balance - toAccount.Balance

		require.Equal(t, diffBalanceFromAccount, diffBalanceToAccount)
		require.True(t, diffBalanceFromAccount > 0)
		require.True(t, diffBalanceFromAccount%amount == 0) // 1 * amount, 2 * amount, ..., n * amount

		k := int(diffBalanceFromAccount / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final update balances
	updatedFromAccount, err := testQueries.GetAccount(context.Background(), fromAccount.ID)
	require.NoError(t, err)

	updatedToAccount, err := testQueries.GetAccount(context.Background(), toAccount.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedFromAccount.Balance, updatedToAccount.Balance)
	require.Equal(t, fromAccount.Balance-int64(n)*amount, updatedFromAccount.Balance)
	require.Equal(t, toAccount.Balance+int64(n)*amount, updatedToAccount.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := db.NewStore(testDB)

	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	fmt.Println(">> before:", fromAccount.Balance, toAccount.Balance)

	n := 10
	amount := int64(10)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := fromAccount.ID
		toAccountID := toAccount.ID

		if i%2 == 1 {
			fromAccountID = toAccount.ID
			toAccountID = fromAccount.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), db.TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check the final updated balance
	updatedAccount1, err := store.GetAccount(context.Background(), fromAccount.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), toAccount.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, fromAccount.Balance, updatedAccount1.Balance)
	require.Equal(t, toAccount.Balance, updatedAccount2.Balance)
}
