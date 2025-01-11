package app

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/charles-m-knox/investment-balancer/pkg/balancer"
	"github.com/charles-m-knox/investment-balancer/pkg/quote"
)

type Config struct {
	Strategies         map[string]balancer.Strategy `json:"strategies"`
	Accounts           map[string]balancer.Account  `json:"accounts"`
	QuoteCache         quote.QuoteCacheMap          `json:"quoteCache"`
	AlphaVantageAPIKey string                       `json:"alphaVantageApiKey"`

	configFilePath string `json:"-"`
}

func (conf *Config) LoadConfig(configFilePath string, appName string, configFileName string) error {
	var err error
	if configFilePath == "" {
		configFilePath, err = xdg.SearchConfigFile(path.Join(appName, configFileName))
		if err != nil {
			return fmt.Errorf("failed to get xdg config dir: %v", err.Error())
		}
	}

	if configFilePath != "" {
		bac, err := os.ReadFile(configFilePath)
		if err != nil {
			return fmt.Errorf("config file not readable at %v", configFilePath)
		}

		err = json.Unmarshal(bac, conf)
		if err != nil {
			return fmt.Errorf("config file %v failed to parse: %v", configFilePath, err.Error())
		}

		log.Printf("loaded config from %v", configFilePath)
	} else {
		if xdg.ConfigHome != "" {
			configFilePath = path.Join(xdg.ConfigHome, appName, configFileName)
			log.Printf("using %v for config file path", configFilePath)
		} else {
			return fmt.Errorf("unable to automatically identify any suitable config dirs; configuration will not be saved")
		}
	}

	conf.configFilePath = configFilePath

	return nil
}

func (conf *Config) SaveConfig() error {
	if conf.configFilePath == "" {
		return fmt.Errorf("received empty config filename")
	}

	if conf == nil {
		return fmt.Errorf("config was nil")
	}

	b, err := json.Marshal(conf)
	if err != nil {
		return fmt.Errorf("failed to marshal app config to json: %v", err.Error())
	} else {
		dir, _ := filepath.Split(conf.configFilePath)
		if dir != "" {
			err := os.MkdirAll(dir, 0o755)
			if err != nil {
				return fmt.Errorf("failed to create app config parent dir %v: %v", dir, err.Error())
			}
		}

		err = os.WriteFile(conf.configFilePath, b, 0o644)
		if err != nil {
			return fmt.Errorf("failed to save app config to %v: %v", conf.configFilePath, err.Error())
		}
	}

	return nil
}
