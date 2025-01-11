package balancer

import (
	"github.com/shopspring/decimal"
)

type Symbol struct {
	Symbol string `json:"symbol"`
	Type   string `json:"type"`
}

type Strategy struct {
	Name        string                     `json:"name"`
	Symbols     []Symbol                   `json:"symbols"`
	Allocations map[string]decimal.Decimal `json:"allocations"`
}

type CachedQuote struct {
	Quote `json:"quote"`
	Time  int64 `json:"time"`
}

type Account struct {
	Name     string          `json:"name"`
	Balance  decimal.Decimal `json:"balance"`
	Strategy string          `json:"strategy"`
}

type Allocation struct {
	Shares          int64           `json:"shares" `
	SharePrice      decimal.Decimal `json:"sharePrice" `
	Remainder       decimal.Decimal `json:"remainder" `
	TotalAllocated  decimal.Decimal `json:"totalAllocated" `
	IdealAllocation decimal.Decimal `json:"idealAllocation" `
	// The ideal allocation out of the total balance divided up for this group.
	IdealGroupAllocationPercentage decimal.Decimal `json:"idealGroupAllocationPercentage" `
	// The ideal allocation out of the grand total balance for the entire
	// account.
	IdealSymbolAllocationPercentage decimal.Decimal `json:"idealSymbolAllocationPercentage" `
}

type Quote struct {
	Symbol string          `json:"symbol" `
	Price  decimal.Decimal `json:"price" `
}

type StockData struct {
	GlobalQuote struct {
		Symbol string `json:"01. symbol"`
		Price  string `json:"05. price"`
	} `json:"Global Quote"`
}
