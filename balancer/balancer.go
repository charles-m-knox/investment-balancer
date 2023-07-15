package balancer

import (
	"fmt"

	"investment-balancer-v3/config"
	"investment-balancer-v3/models"

	"github.com/shopspring/decimal"
)

func BalanceAccount(conf config.Config, bal models.Account, quotes []models.Quote) (result map[string]map[string]models.Allocation, err error) {
	// divide a balance according to the portfolio
	portfolio, err := conf.GetPortfolio(bal.Strategy)
	if err != nil {
		return result, fmt.Errorf(
			"failed to balance: %v",
			err.Error(),
		)
	}

	// first, group symbols according to their classification
	groups := make(map[string]map[string]models.Allocation)
	for _, symbol := range portfolio.Symbols {
		if groups[symbol.Type] == nil {
			groups[symbol.Type] = make(map[string]models.Allocation)
		}
		groups[symbol.Type][symbol.Symbol] = models.Allocation{}
	}

	// now that we have all symbols grouped, proceed
	// to apply the allocations
	for group, symbols := range groups {
		groupAllocation := portfolio.Allocations[group].Div(
			decimal.NewFromInt(100),
		).Mul(
			bal.Balance,
		)

		numSymbols := int64(len(symbols))
		numSymbolsDec := decimal.NewFromInt(numSymbols)

		allocPercentageFromTotal := portfolio.Allocations[group].Div(
			numSymbolsDec,
		)

		allocPerSymbol := groupAllocation.Div(
			numSymbolsDec,
		)

		// the balance for this type of investment has been established,
		// so proceed to skim over each symbol associated with this type
		// of investment and find out how many shares to buy
		for symbol := range symbols {
			for _, quote := range quotes {
				if quote.Symbol != symbol {
					continue
				}

				shares := allocPerSymbol.Div(quote.Price).Floor()
				totalAllocated := shares.Mul(quote.Price)

				groups[group][symbol] = models.Allocation{
					Shares:                          shares.IntPart(), // will always be accurate due to earlier Floor()
					SharePrice:                      quote.Price,
					Remainder:                       allocPerSymbol.Sub(totalAllocated),
					TotalAllocated:                  totalAllocated,
					IdealAllocation:                 allocPerSymbol,
					IdealGroupAllocationPercentage:  portfolio.Allocations[group],
					IdealSymbolAllocationPercentage: allocPercentageFromTotal,
				}
			}
		}
	}

	return groups, nil
}
