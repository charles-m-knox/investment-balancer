package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"

	"investment-balancer-v3/models"

	"gopkg.in/yaml.v3"
)

const (
	ConfigFile = "config.yml"
)

type Config struct {
	OutputFilename     string            `yaml:"outputFilename"`
	Strategies         []models.Strategy `yaml:"strategies"`
	Accounts           []models.Account  `yaml:"accounts"`
	AlphaVantageAPIKey string            `yaml:"alphaVantageApiKey"`
}

func (c *Config) AssertValidConfig() {
	if c.AlphaVantageAPIKey == "" {
		log.Fatalf("apiKey is not set")
	}
}

// LoadConfig reads from a provided yaml-formatted configuration filename
func LoadConfig() (conf Config, err error) {
	f, err := os.Open("config.yml")
	if err != nil {
		return conf, fmt.Errorf("failed to open config.yml: %v", err)
	}

	c, err := io.ReadAll(f)
	if err != nil {
		return conf, fmt.Errorf("failed to ReadAll from file: %v", err)
	}

	err = yaml.Unmarshal(c, &conf)
	if err != nil {
		return conf, fmt.Errorf("failed to unmarshal yaml config: %v", err)
	}

	return conf, nil
}

// GetAllSymbols uses a map to retrieve all unique ticker symbols
// across all portfolios from the config, and then returns them
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

// GetPortfolio attempts to retrieve a portfolio by name. If it cannot find one by the provided name, it will return an error.
func (conf *Config) GetPortfolio(name string) (result models.Strategy, err error) {
	for _, portfolio := range conf.Strategies {
		if portfolio.Name == name {
			return portfolio, nil
		}
	}

	return result, fmt.Errorf(
		"failed to find a portfolio by name %v",
		name,
	)
}
