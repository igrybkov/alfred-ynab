package ynab

import (
	api "github.com/brunomvsouza/ynab.go"
	categoryApi "github.com/brunomvsouza/ynab.go/api/category"
)

type Category = categoryApi.Category
type Group = categoryApi.Group
type GroupWithCategories = categoryApi.GroupWithCategories

func GetCategories(client api.ClientServicer, budgetID string) ([]*GroupWithCategories, error) {
	categories, err := client.Category().GetCategories(budgetID, nil)
	if err != nil {
		return nil, err
	}
	return categories.GroupWithCategories, nil
}
