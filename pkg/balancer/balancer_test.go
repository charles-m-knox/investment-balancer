package balancer

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
)

func Test_BalanceAccount(t *testing.T) {
	dollars0 := decimal.NewFromInt(0)
	percent25 := decimal.NewFromInt(25)
	dollars25 := decimal.NewFromInt(25)
	dollars50 := decimal.NewFromInt(50)
	percent50 := decimal.NewFromInt(50)
	dollars100 := decimal.NewFromInt(100)

	testStrategy := "test-strategy"
	acct100 := Account{Name: "100dollars", Balance: dollars100, Strategy: testStrategy}

	tests := []struct {
		name     string
		expected map[string]map[string]Allocation
		err      error
		conf     Config
		bal      Account
		quotes   []Quote
	}{
		{
			name: "balances a simple strategy",
			err:  nil,
			conf: Config{
				Strategies: []Strategy{
					{
						Name: testStrategy,
						Symbols: []Symbol{
							{"SCHX", "large"},
							{"SCHB", "large"},
							{"SCHD", "dividends"},
						},
						Allocations: map[string]decimal.Decimal{
							"large":     percent50,
							"dividends": percent50,
						},
					},
				},
				Accounts: []Account{acct100},
			},
			quotes: []Quote{
				{Symbol: "SCHD", Price: dollars25},
				{Symbol: "SCHB", Price: dollars25},
				{Symbol: "SCHX", Price: dollars25},
			},
			bal: acct100,
			expected: map[string]map[string]Allocation{
				"large": {
					"SCHX": {
						Shares:                          1,
						SharePrice:                      dollars25,
						Remainder:                       dollars0,
						TotalAllocated:                  dollars25,
						IdealAllocation:                 dollars25,
						IdealGroupAllocationPercentage:  percent50,
						IdealSymbolAllocationPercentage: percent25,
					},
					"SCHB": {
						Shares:                          1,
						SharePrice:                      dollars25,
						Remainder:                       dollars0,
						TotalAllocated:                  dollars25,
						IdealAllocation:                 dollars25,
						IdealGroupAllocationPercentage:  percent50,
						IdealSymbolAllocationPercentage: percent25,
					},
				},
				"dividends": {
					"SCHD": {
						Shares:                          2,
						SharePrice:                      dollars25,
						Remainder:                       dollars0,
						TotalAllocated:                  dollars50,
						IdealAllocation:                 dollars50,
						IdealGroupAllocationPercentage:  percent50,
						IdealSymbolAllocationPercentage: percent50,
					},
				},
			},
		},
		{
			name: "works with underscores in symbol names",
			err:  nil,
			conf: Config{
				Strategies: []Strategy{
					{
						Name: testStrategy,
						Symbols: []Symbol{
							{"SCHD__1", "large"},
							{"SCHB", "large"},
							{"SCHD__2", "dividends"},
						},
						Allocations: map[string]decimal.Decimal{
							"large":     percent50,
							"dividends": percent50,
						},
					},
				},
				Accounts: []Account{acct100},
			},
			quotes: []Quote{
				{Symbol: "SCHD", Price: dollars25},
				{Symbol: "SCHB", Price: dollars25},
			},
			bal: acct100,
			expected: map[string]map[string]Allocation{
				"large": {
					"SCHD__1": {
						Shares:                          1,
						SharePrice:                      dollars25,
						Remainder:                       dollars0,
						TotalAllocated:                  dollars25,
						IdealAllocation:                 dollars25,
						IdealGroupAllocationPercentage:  percent50,
						IdealSymbolAllocationPercentage: percent25,
					},
					"SCHB": {
						Shares:                          1,
						SharePrice:                      dollars25,
						Remainder:                       dollars0,
						TotalAllocated:                  dollars25,
						IdealAllocation:                 dollars25,
						IdealGroupAllocationPercentage:  percent50,
						IdealSymbolAllocationPercentage: percent25,
					},
				},
				"dividends": {
					"SCHD__2": {
						Shares:                          2,
						SharePrice:                      dollars25,
						Remainder:                       dollars0,
						TotalAllocated:                  dollars50,
						IdealAllocation:                 dollars50,
						IdealGroupAllocationPercentage:  percent50,
						IdealSymbolAllocationPercentage: percent50,
					},
				},
			},
		},
		{
			name: "throws an error if the strategy isn't present",
			err:  fmt.Errorf("failed to balance"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := BalanceAccount(test.conf, test.bal, test.quotes)

			if test.err == nil && err != nil {
				t.Fatalf("test.err is nil but actual err occurred: %v", err.Error())
			} else if test.err != nil && err == nil {
				t.Fatalf("test.err expected an error but did not get one")
			}

			for expectedKey, expectedVal := range test.expected {
				found := false
				for actualKey, actualVal := range actual {
					if actualKey == expectedKey {
						found = true
						for expectedAllocationKey, expectedAllocationVal := range expectedVal {
							foundAlloc := false
							for actualAllocationKey, actualAllocationVal := range actualVal {
								if expectedAllocationKey == actualAllocationKey {
									foundAlloc = true
									if !expectedAllocationVal.IdealAllocation.Equal(actualAllocationVal.IdealAllocation) {
										t.Errorf("test.expected[%v][%v].IdealAllocation = %v, but actual[%v][%v].IdealAllocation = %v", expectedKey, expectedAllocationKey, expectedAllocationVal.IdealAllocation, actualKey, actualAllocationKey, actualAllocationVal.IdealAllocation)
									}
									if !expectedAllocationVal.IdealGroupAllocationPercentage.Equal(actualAllocationVal.IdealGroupAllocationPercentage) {
										t.Errorf("test.expected[%v][%v].IdealGroupAllocationPercentage = %v, but actual[%v][%v].IdealGroupAllocationPercentage = %v", expectedKey, expectedAllocationKey, expectedAllocationVal.IdealGroupAllocationPercentage, actualKey, actualAllocationKey, actualAllocationVal.IdealAllocation)
									}
									if !expectedAllocationVal.IdealSymbolAllocationPercentage.Equal(actualAllocationVal.IdealSymbolAllocationPercentage) {
										t.Errorf("test.expected[%v][%v].IdealSymbolAllocationPercentage = %v, but actual[%v][%v].IdealSymbolAllocationPercentage = %v", expectedKey, expectedAllocationKey, expectedAllocationVal.IdealSymbolAllocationPercentage, actualKey, actualAllocationKey, actualAllocationVal.IdealAllocation)
									}
									if !expectedAllocationVal.Remainder.Equal(actualAllocationVal.Remainder) {
										t.Errorf("test.expected[%v][%v].Remainder = %v, but actual[%v][%v].Remainder = %v", expectedKey, expectedAllocationKey, expectedAllocationVal.Remainder, actualKey, actualAllocationKey, actualAllocationVal.IdealAllocation)
									}
									if !expectedAllocationVal.SharePrice.Equal(actualAllocationVal.SharePrice) {
										t.Errorf("test.expected[%v][%v].SharePrice = %v, but actual[%v][%v].SharePrice = %v", expectedKey, expectedAllocationKey, expectedAllocationVal.SharePrice, actualKey, actualAllocationKey, actualAllocationVal.IdealAllocation)
									}
									if expectedAllocationVal.Shares != actualAllocationVal.Shares {
										t.Errorf("test.expected[%v][%v].Shares = %v, but actual[%v][%v].Shares = %v", expectedKey, expectedAllocationKey, expectedAllocationVal.Shares, actualKey, actualAllocationKey, actualAllocationVal.IdealAllocation)
									}
									if !expectedAllocationVal.TotalAllocated.Equal(actualAllocationVal.TotalAllocated) {
										t.Errorf("test.expected[%v][%v].TotalAllocated = %v, but actual[%v][%v].TotalAllocated = %v", expectedKey, expectedAllocationKey, expectedAllocationVal.TotalAllocated, actualKey, actualAllocationKey, actualAllocationVal.IdealAllocation)
									}
								}
							}

							if !foundAlloc {
								t.Errorf("test.expected[%v][%v] not found in actual", expectedKey, expectedAllocationKey)
							}
						}
					}
				}
				if !found {
					t.Errorf("test.expected[%v] not found in actual", expectedKey)
				}
			}
		})
	}
}
