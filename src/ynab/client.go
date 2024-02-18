package ynab

import api "github.com/brunomvsouza/ynab.go"

func MakeYNABClient(accessToken string) api.ClientServicer {
	return api.NewClient(accessToken)
}
