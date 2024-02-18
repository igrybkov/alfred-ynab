package main

import (
	"encoding/json"
	"log"
	"strconv"
	"sync"
	"time"

	api "github.com/brunomvsouza/ynab.go"
	aw "github.com/deanishe/awgo"

	"com.grybkov.alfred-ynab/src/config"
	"com.grybkov.alfred-ynab/src/ynab"
)

type ItemHandler interface {
	Preload(wg *sync.WaitGroup, apiClient api.ClientServicer, cacheTTL time.Duration)
	AddItems(workflow *aw.Workflow)
}

type TransactionsHandler struct {
	unapprovedTransactions    []*ynab.Transaction
	uncategorizedTransactions []*ynab.Transaction
}

func (t *TransactionsHandler) Preload(wg *sync.WaitGroup, apiClient api.ClientServicer, cacheTTL time.Duration) {
	go func() {
		defer wg.Done()
		wg.Add(1)
		go func() {
			reloadUnapprovedTransactions := func() (interface{}, error) {
				return ynab.GetUnapprovedTransactions(apiClient, config.GetYnabBudget())
			}
			defer wg.Done()
			err := wf.Cache.LoadOrStoreJSON("unapprovedTransactions", cacheTTL, reloadUnapprovedTransactions, &t.unapprovedTransactions)
			if err != nil {
				return
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			reloadUncategorizedTransactions := func() (interface{}, error) {
				return ynab.GetUncategorizedTransactions(apiClient, config.GetYnabBudget())
			}
			err := wf.Cache.LoadOrStoreJSON("uncategorizedTransactions", cacheTTL, reloadUncategorizedTransactions, &t.uncategorizedTransactions)
			if err != nil {
				return
			}
		}()
	}()
}

func (t *TransactionsHandler) AddItems(workflow *aw.Workflow) {
	icon := aw.Icon{Value: "icons/hand-coins.svg"}
	totalUnapprovedTransactions := strconv.Itoa(len(t.unapprovedTransactions))
	unapprovedArgJson, _ := json.Marshal(Arg{EntityID: "", EntityType: "transaction", Value: totalUnapprovedTransactions, BrowserURL: ynab.BuildBudgetURL(config.GetYnabBudget())})
	workflow.NewItem("Unapproved transactions: " + totalUnapprovedTransactions).
		Arg(string(unapprovedArgJson)).
		Valid(true).
		Icon(&icon)

	totalUncategorizedTransactions := strconv.Itoa(len(t.uncategorizedTransactions))
	uncategorizedArgJson, _ := json.Marshal(Arg{EntityID: "", EntityType: "transaction", Value: totalUncategorizedTransactions, BrowserURL: ynab.BuildBudgetURL(config.GetYnabBudget())})
	workflow.NewItem("Uncategorized Transactions: " + totalUncategorizedTransactions).
		Arg(string(uncategorizedArgJson)).
		Valid(true).
		Icon(&icon)

}

type AccountsHandler struct {
	accounts []*ynab.Account
}

func (a *AccountsHandler) Preload(wg *sync.WaitGroup, apiClient api.ClientServicer, cacheTTL time.Duration) {
	log.Println("Preloading accounts")
	go func() {
		defer wg.Done()
		reloadAccounts := func() (interface{}, error) {
			log.Println("Reloading accounts")
			return ynab.GetAccounts(apiClient, config.GetYnabBudget())
		}
		err := wf.Cache.LoadOrStoreJSON("accounts", cacheTTL, reloadAccounts, &a.accounts)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Accounts reloaded")
	}()
}

func (a *AccountsHandler) AddItems(workflow *aw.Workflow) {
	log.Println("Adding account items")
	var iconPerAccountType = map[ynab.AccountType]aw.Icon{
		"creditCard": aw.Icon{Value: "icons/credit-card.svg"},
		"default":    aw.Icon{Value: "icons/bank.svg"},
	}
	for _, account := range a.accounts {
		if account.Deleted || account.Closed {
			continue
		}

		balance := strconv.FormatFloat(float64(account.Balance)/1000, 'g', -1, 64)
		argJson, _ := json.Marshal(Arg{EntityID: account.ID, EntityType: "account", Value: balance, BrowserURL: ynab.BuildAccountURL(config.GetYnabBudget(), account.ID)})
		item := wf.NewItem(account.Name + ": " + balance).
			Subtitle(account.Name).
			Arg(string(argJson)).
			Valid(true)
		if icon, ok := iconPerAccountType[account.Type]; ok {
			item.Icon(&icon)
		} else {
			icon := iconPerAccountType["default"]
			item.Icon(&icon)
		}
	}
	log.Println("Account items added")
}

type CategoriesHandler struct {
	categories []*ynab.GroupWithCategories
}

func (c *CategoriesHandler) Preload(wg *sync.WaitGroup, apiClient api.ClientServicer, cacheTTL time.Duration) {
	go func() {
		defer wg.Done()
		reloadCategories := func() (interface{}, error) {
			return ynab.GetCategories(apiClient, config.GetYnabBudget())
		}
		err := wf.Cache.LoadOrStoreJSON("categories", cacheTTL, reloadCategories, &c.categories)
		if err != nil {
			return
		}
	}()
}

func (c *CategoriesHandler) AddItems(workflow *aw.Workflow) {
	for _, group := range c.categories {
		if group.Deleted || group.Hidden {
			continue
		}
		for _, category := range group.Categories {
			if category.Deleted || category.Hidden {
				continue
			}
			balance := strconv.FormatFloat(float64(category.Balance)/1000, 'g', -1, 64)
			item := wf.NewItem(category.Name + ": " + balance)
			icon := aw.Icon{Value: "icons/piggy-bank.svg"}
			argJson, _ := json.Marshal(Arg{EntityID: category.ID, EntityType: "category", Value: balance, BrowserURL: ynab.BuildBudgetURL(config.GetYnabBudget())})
			item.
				Subtitle(group.Name).
				Arg(string(argJson)).
				Icon(&icon).
				Valid(true)
		}
	}
}
