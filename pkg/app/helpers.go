package app

import (
	"sort"

	"github.com/charles-m-knox/investment-balancer/pkg/balancer"
)

// GetAllSymbols uses a map to retrieve all unique ticker symbols across all
// portfolios from the config, and then returns them
func (conf *Config) GetAllSymbols() (symbols []string) {
	uniqueSymbols := make(map[string]string)

	for _, portfolio := range conf.Strategies {
		for _, symbol := range portfolio.Symbols {
			uniqueSymbols[symbol.Symbol] = symbol.Symbol
		}
	}

	for symbol := range uniqueSymbols {
		symbols = append(symbols, symbol)
	}

	sort.Strings(symbols)

	return
}

func (conf *Config) GetFirstAccountKey() string {
	for key := range conf.Accounts {
		return key
	}
	return ""
}

func (conf *Config) GetFirstAccount() (balancer.Account, bool) {
	acct, ok := conf.Accounts[conf.GetFirstAccountKey()]
	return acct, ok
}
