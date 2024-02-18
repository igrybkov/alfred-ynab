package ynab

import (
	api "github.com/brunomvsouza/ynab.go"
	accountApi "github.com/brunomvsouza/ynab.go/api/account"
)

type Account = accountApi.Account
type AccountType = accountApi.Type

func GetAccounts(client api.ClientServicer, budgetID string) ([]*Account, error) {
	accounts, err := client.Account().GetAccounts(budgetID, nil)
	if err != nil {
		return nil, err
	}

	return accounts.Accounts, nil
}
