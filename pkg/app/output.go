package app

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/charles-m-knox/investment-balancer/pkg/balancer"
	"github.com/shopspring/decimal"
)

func (conf *Config) PrettyPrint(quotes []balancer.Quote) (string, error) {
	headers := []string{}

	if len(conf.Accounts) > 1 {
		headers = append(headers, "Name")
	}

	headers = append(headers,
		"Symbol",
		"Type",
		"Shares",
		"Share_Price",
		"Purchase_Price",
		"Allocated",
		"Remainder",
		"Symbol_Allocation_%",
		"Group_Allocation_%",
		"From_Balance",
	)
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	w.Comma = '\t'

	err := w.Write(headers)
	if err != nil {
		return "", fmt.Errorf("error writing record to csv: %v", err.Error())
	}

	for _, account := range conf.Accounts {
		strategy, ok := conf.Strategies[account.Strategy]
		if !ok {
			return "", fmt.Errorf("failed to find strategy: %v", account.Strategy)
		}

		groups, err := balancer.BalanceAccount(strategy, account, quotes)
		if err != nil {
			return "", fmt.Errorf("failed to balance: %v", err.Error())
		}

		for group, symbols := range groups {
			for symbol, s := range symbols {
				r := []string{}
				if len(conf.Accounts) > 1 {
					r = append(r, account.Name)
				}

				r = append(r,
					symbol,
					group,
					fmt.Sprintf("%v", s.Shares),
					fmt.Sprintf("%v", s.SharePrice.Truncate(2)),
					fmt.Sprintf("%v", s.TotalAllocated.Truncate(2)),
					fmt.Sprintf("%v", s.IdealAllocation.Truncate(2)),
					fmt.Sprintf("%v", s.Remainder.Truncate(2)),
					fmt.Sprintf("%v", s.IdealSymbolAllocationPercentage.Div(decimal.NewFromInt(100)).Truncate(2)),
					fmt.Sprintf("%v", s.IdealGroupAllocationPercentage.Div(decimal.NewFromInt(100)).Truncate(2)),
					fmt.Sprintf("%v", account.Balance.Truncate(2)),
				)
				err := w.Write(r)
				if err != nil {
					return "", fmt.Errorf("failed to write newRecord: %v", err.Error())
				}
			}
		}
	}

	w.Flush()

	err = w.Error()
	if err != nil {
		return "", fmt.Errorf("failed to write csv: %v", err.Error())
	}

	allData := buf.String()
	lines := strings.Split(allData, "\n")

	var buf2 bytes.Buffer
	wt := tabwriter.NewWriter(&buf2, 0, 0, 2, ' ', 0)

	for _, line := range lines {
		fmt.Fprintln(wt, line)
	}

	wt.Flush()

	return fmt.Sprintf("%v", buf2.String()), nil
}
