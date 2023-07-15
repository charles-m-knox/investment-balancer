package models

import (
	"github.com/shopspring/decimal"
)

type Symbol struct {
	Symbol string `yaml:"symbol"`
	Type   string `yaml:"type"`
}

type Strategy struct {
	Name        string                     `yaml:"name"`
	Symbols     []Symbol                   `yaml:"symbols"`
	Allocations map[string]decimal.Decimal `yaml:"allocations"`
}

type CachedQuote struct {
	Quote `json:"quote"`
	Time  int64 `json:"time"`
}

type Account struct {
	Name     string          `yaml:"name"`
	Balance  decimal.Decimal `yaml:"balance"`
	Strategy string          `yaml:"strategy"`
}

type Allocation struct {
	Shares                          int64
	SharePrice                      decimal.Decimal
	Remainder                       decimal.Decimal
	TotalAllocated                  decimal.Decimal
	IdealAllocation                 decimal.Decimal
	IdealGroupAllocationPercentage  decimal.Decimal
	IdealSymbolAllocationPercentage decimal.Decimal
}

type Quote struct {
	Symbol string
	Price  decimal.Decimal
}

type StockData struct {
	GlobalQuote struct {
		Symbol string `json:"01. symbol"`
		Price  string `json:"05. price"`
	} `json:"Global Quote"`
}
