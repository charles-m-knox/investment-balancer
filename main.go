package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"investment-balancer-v3/balancer"
	"investment-balancer-v3/config"
	"investment-balancer-v3/helpers"
	"investment-balancer-v3/models"

	d "github.com/shopspring/decimal"
)

var QuoteCache map[string]models.CachedQuote

func loadCache() error {
	_, err := os.Stat(".quotecache.json")
	if errors.Is(err, fs.ErrNotExist) {
		log.Println("no symbol cache exists yet")
		return nil
	}

	f, err := os.Open(".quotecache.json")
	if err != nil {
		return fmt.Errorf("failed to load symbol cache: %v", err)
	}

	c, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("failed to read symbol cache: %v", err)
	}

	err = json.Unmarshal(c, &QuoteCache)
	if err != nil {
		return fmt.Errorf("failed to unmarshal symbol cache: %v", err)
	}

	return nil
}

func saveCache() error {
	b, err := json.Marshal(QuoteCache)
	if err != nil {
		return fmt.Errorf("failed to marshal symbol cache: %v", err)
	}

	err = os.WriteFile(".quotecache.json", b, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write symbol cache: %v", err)
	}

	return nil
}

// getLatestPrice  returns the symbol's latest price, a boolean indicating if
// the result was retrieved from cache (to assist with rate limiting), and
// an error, if encountered
func getLatestPrice(symbol string, apiKey string) (d.Decimal, bool, error) {
	// special case: all cash allocations are instant
	if symbol == "_cash" {
		QuoteCache[symbol] = models.CachedQuote{
			Quote: models.Quote{
				Symbol: "_cash",
				Price:  d.NewFromInt(1),
			},
			Time: time.Now().Unix(),
		}
	}

	cached, ok := QuoteCache[symbol]
	// cached symbols are good for 6 hours (I chose this randomly)
	now := time.Now()
	nowUnix := now.Unix()
	isCached := helpers.IsWithin(cached.Time, nowUnix, 6*time.Hour)
	if ok && isCached {
		log.Printf("cached price for %v", symbol)
		return cached.Price, true, nil
	}

	apiURL := fmt.Sprintf(
		"https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=%s&apikey=%s",
		symbol,
		apiKey,
	)

	resp, err := http.Get(apiURL)
	if err != nil {
		return d.Zero, false, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return d.Zero, false, err
	}

	log.Println(string(body))

	var data models.StockData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return d.Zero, false, err
	}

	price, err := d.NewFromString(data.GlobalQuote.Price)
	if err != nil {
		log.Fatalf(
			"failed to parse decimal from string %v: %v", price, err,
		)
	}

	log.Printf("caching quote for symbol %v...", symbol)
	QuoteCache[symbol] = models.CachedQuote{
		Quote: models.Quote{
			Symbol: symbol,
			Price:  price,
		},
		Time: time.Now().Unix(),
	}

	return price, false, nil
}

func main() {
	QuoteCache = make(map[string]models.CachedQuote)

	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err.Error())
	}

	f, err := os.Create(conf.OutputFilename)
	if err != nil {
		log.Fatalf(
			"failed to open %v for writing: %v", conf.OutputFilename, err,
		)
	}

	conf.AssertValidConfig()
	log.Println("config is valid")

	symbols := conf.GetAllSymbols()
	quotes := []models.Quote{}

	err = loadCache()
	if err != nil {
		log.Fatalf("symbol cache load failure: %v", err)
	}

	for _, symbol := range symbols {
		log.Printf("attempting to get price for symbol %v...", symbol)
		price, cached, err := getLatestPrice(symbol, conf.AlphaVantageAPIKey)
		if err != nil {
			log.Printf("error fetching price for %v: %v\n", symbol, err)
		}

		log.Printf("The latest price for %v is: %v\n", symbol, price)

		quotes = append(quotes, models.Quote{
			Symbol: symbol,
			Price:  price,
		})

		if cached { // no need to rate limit if using cache
			continue
		}

		err = saveCache()
		if err != nil {
			log.Fatalf("symbol cache save failure: %v", err)
		}

		time.Sleep(1 * time.Second)
	}

	err = saveCache()
	if err != nil {
		log.Fatalf("symbol cache save failure: %v", err)
	}

	// write CSV headers
	headers := []string{
		"Name",
		"Symbol",
		"Type",
		"Shares",
		"Share Price",
		"Purchase Price",
		"Allocated",
		"Remainder",
		"Symbol Allocation %",
		"Group Allocation %",
		"From Balance",
	}

	w := csv.NewWriter(f)

	err = w.Write(headers)
	if err != nil {
		log.Fatalln("error writing record to csv:", err)
	}

	for _, balance := range conf.Accounts {
		// for each balance in conf.Balances, proceed to apply the portfolio
		groups, err := balancer.BalanceAccount(conf, balance, quotes)
		if err != nil {
			log.Fatalf("failed to balance: %v", err.Error())
		}

		for group, symbols := range groups {
			for symbol, s := range symbols {
				err := w.Write([]string{
					balance.Name,                // "Name"
					symbol,                      // "Symbol"
					group,                       // "Type"
					fmt.Sprintf("%v", s.Shares), // "Shares"
					fmt.Sprintf("%v", s.SharePrice.Truncate(2)),                                             // $ "Shares"
					fmt.Sprintf("%v", s.TotalAllocated.Truncate(2)),                                         // $ "Purchase Price"
					fmt.Sprintf("%v", s.IdealAllocation.Truncate(2)),                                        // $ "Allocated"
					fmt.Sprintf("%v", s.Remainder.Truncate(2)),                                              // $ "Remainder"
					fmt.Sprintf("%v", s.IdealSymbolAllocationPercentage.Div(d.NewFromInt(100)).Truncate(2)), // "Allocation %"
					fmt.Sprintf("%v", s.IdealGroupAllocationPercentage.Div(d.NewFromInt(100)).Truncate(2)),  // "Allocation %"
					fmt.Sprintf("%v", balance.Balance.Truncate(2)),                                          // $ "From Balance"
				})
				if err != nil {
					log.Fatalf("failed to write newRecord: %v", err.Error())
				}
			}
		}
	}

	w.Flush()

	err = w.Error()
	if err != nil {
		log.Fatalf("failed to write csv: %v", err.Error())
	}

	log.Printf("done: finished writing to %v", conf.OutputFilename)
}
