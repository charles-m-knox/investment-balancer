package quote

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/charles-m-knox/investment-balancer/pkg/balancer"

	d "github.com/shopspring/decimal"
)

// GetLatestPrice returns the symbol's latest price, as well as a boolean
// indicating if the result was retrieved from cache (to assist with rate
// limiting).
//
// You do not need to manually update your quote cache, this function will do
// it for you.
func GetLatestPrice(symbol string, apiKey string, quoteCache QuoteCacheMap) (d.Decimal, bool, error) {
	// special case: all cash allocations are instant
	if symbol == "_cash" {
		quoteCache[symbol] = balancer.CachedQuote{
			Quote: balancer.Quote{
				Symbol: "_cash",
				Price:  d.NewFromInt(1),
			},
			Time: time.Now().Unix(),
		}
	}

	// supports duplicate symbols but in different allocations (i.e. SCHD__1,
	// SCHD__2)
	symbol = strings.Split(symbol, "__")[0]

	cached, ok := quoteCache[symbol]

	// cached symbols are good for 6 hours (I chose this randomly)
	now := time.Now()
	nowUnix := now.Unix()

	isCached := balancer.IsWithin(cached.Time, nowUnix, 6*time.Hour)
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return d.Zero, false, err
	}

	log.Println(string(body))

	var data balancer.StockData
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

	quoteCache[symbol] = balancer.CachedQuote{
		Quote: balancer.Quote{
			Symbol: symbol,
			Price:  price,
		},
		Time: time.Now().Unix(),
	}

	return price, false, nil
}
