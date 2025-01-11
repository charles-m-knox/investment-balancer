package app

import (
	"log"
	"time"

	"github.com/charles-m-knox/investment-balancer/pkg/balancer"
	"github.com/charles-m-knox/investment-balancer/pkg/quote"
)

func (conf *Config) GetQuotes() []balancer.Quote {
	quotes := []balancer.Quote{}
	for _, symbol := range conf.GetAllSymbols() {
		log.Printf("attempting to get price for symbol %v...", symbol)
		price, cached, err := quote.GetLatestPrice(symbol, conf.AlphaVantageAPIKey, conf.QuoteCache)
		if err != nil {
			log.Printf("error fetching price for %v: %v\n", symbol, err)
		}

		log.Printf("The latest price for %v is: %v\n", symbol, price)

		quotes = append(quotes, balancer.Quote{
			Symbol: symbol,
			Price:  price,
		})

		if cached { // no need to rate limit if using cache
			continue
		}

		err = conf.SaveConfig()
		if err != nil {
			log.Printf("failed to save cache to config: %v", err.Error())
		}

		time.Sleep(1 * time.Second)
	}

	return quotes
}
