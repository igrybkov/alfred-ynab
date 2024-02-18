package ynab

const baseURL = "https://app.ynab.com/"

func BuildURL(budgetID string, path string) string {
	return baseURL + budgetID + "/" + path
}

func BuildBudgetURL(budgetID string) string {
	return BuildURL(budgetID, "budget")
}

func BuildAccountURL(budgetID string, accountID string) string {
	return BuildURL(budgetID, "accounts/"+accountID)
}
