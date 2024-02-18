package ynab

import (
	api "github.com/brunomvsouza/ynab.go"
	transactionApi "github.com/brunomvsouza/ynab.go/api/transaction"
)

type Transaction = transactionApi.Transaction

func filterTransactions(transactions []*Transaction, filter func(*Transaction) bool) []*Transaction {
	var result = make([]*Transaction, 0)
	for _, transaction := range transactions {
		if filter(transaction) {
			result = append(result, transaction)
		}
	}
	return result

}

func GetUnapprovedTransactions(client api.ClientServicer, budgetID string) ([]*Transaction, error) {
	var transactions, err = client.Transaction().GetTransactions(budgetID, &transactionApi.Filter{
		Type: transactionApi.StatusUnapproved.Pointer(),
	})
	if err != nil {
		return nil, err
	}

	return filterTransactions(transactions, func(transaction *Transaction) bool {
		return !transaction.Deleted
	}), nil
}

func GetUncategorizedTransactions(client api.ClientServicer, budgetID string) ([]*Transaction, error) {
	var transactions, err = client.Transaction().GetTransactions(budgetID, &transactionApi.Filter{
		Type: transactionApi.StatusUnapproved.Pointer(),
	})
	if err != nil {
		return nil, err
	}

	return filterTransactions(transactions, func(transaction *Transaction) bool {
		return !transaction.Deleted
	}), nil
}
