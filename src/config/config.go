package config

import (
	"strconv"
	"time"

	aw "github.com/deanishe/awgo"
)

var cfg = aw.NewConfig()

func GetYnabAccessToken() string {
	return cfg.Get("YNAB_ACCESS_TOKEN")
}

func GetYnabBudget() string {
	return cfg.Get("YNAB_BUDGET")
}

func GetCacheTTL() time.Duration {
	ttl := cfg.Get("CACHE_TTL", "")
	if ttl == "" {
		return 1 * time.Second
	}
	ttlInt, err := strconv.Atoi(ttl)
	if err != nil {
		return 1 * time.Second
	}
	return time.Duration(ttlInt) * time.Minute
}
