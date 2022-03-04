package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/andrelsf/go-restapi/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	params := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, params.Owner, account.Owner)
	require.Equal(t, params.Balance, account.Balance)
	require.Equal(t, params.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	accountSaved := createRandomAccount(t)

	accountFound, err := testQueries.GetAccount(context.Background(), accountSaved.ID)

	require.NoError(t, err)
	require.NotEmpty(t, accountFound)

	require.Equal(t, accountSaved.ID, accountFound.ID)
	require.Equal(t, accountSaved.Owner, accountFound.Owner)
	require.Equal(t, accountSaved.Balance, accountFound.Balance)
	require.Equal(t, accountSaved.Currency, accountFound.Currency)
	require.WithinDuration(t, accountSaved.CreatedAt, accountFound.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	accountSaved := createRandomAccount(t)

	args := UpdateAccountParams{
		ID:      accountSaved.ID,
		Balance: util.RandomMoney(),
	}

	accountFound, err := testQueries.UpdateAccount(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, accountSaved)

	require.Equal(t, accountSaved.ID, accountFound.ID)
	require.Equal(t, accountSaved.Owner, accountFound.Owner)
	require.Equal(t, args.Balance, accountFound.Balance)
	require.Equal(t, accountSaved.Currency, accountFound.Currency)
	require.WithinDuration(t, accountSaved.CreatedAt, accountFound.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	accountSaved := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), accountSaved.ID)
	require.NoError(t, err)

	accountFound, err := testQueries.GetAccount(context.Background(), accountSaved.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, accountFound)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	params := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), params)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
