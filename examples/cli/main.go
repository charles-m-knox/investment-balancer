package main

import (
	"fmt"
	"log"

	"github.com/charles-m-knox/investment-balancer/pkg/app"
)

const (
	CONFIG_FILE_NAME = "config.json"
	APP_NAME         = "investment-balancer-cli"
)

func main() {
	var conf app.Config

	err := conf.LoadConfig(CONFIG_FILE_NAME, APP_NAME, CONFIG_FILE_NAME)
	if err != nil {
		log.Fatalf("failed to load config: %v", err.Error())
	}

	if conf.AlphaVantageAPIKey == "" {
		log.Fatalf("AlphaVantage API key is not set")
	}

	log.Println("config is valid")

	quotes := conf.GetQuotes()

	output, err := conf.PrettyPrint(quotes)
	if err != nil {
		log.Fatalf("failed to pretty print: %v", err.Error())
	}

	fmt.Printf("%v", output)

	err = conf.SaveConfig()
	if err != nil {
		log.Printf("failed to save cache to config: %v", err.Error())
	}
}
