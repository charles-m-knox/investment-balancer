package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/charles-m-knox/investment-balancer/pkg/balancer"
	"github.com/shopspring/decimal"
)

func (app *App) loadConfig() {
	var err error
	if app.configFilePath == "" {
		app.configFilePath, err = xdg.SearchConfigFile(path.Join(APP_NAME, CONFIG_FILE))
		if err != nil {
			log.Printf("failed to get xdg config dir: %v", err.Error())
		}
	}

	if app.configFilePath != "" {
		bac, err := os.ReadFile(app.configFilePath)
		if err != nil {
			log.Printf("config file not readable at %v", app.configFilePath)
		}

		err = json.Unmarshal(bac, app.conf)
		if err != nil {
			log.Printf("config file %v failed to parse: %v", app.configFilePath, err.Error())
		}

		log.Printf("loaded config from %v", app.configFilePath)
	} else {
		if xdg.ConfigHome != "" {
			app.configFilePath = path.Join(xdg.ConfigHome, APP_NAME, CONFIG_FILE)
			log.Printf("using %v for config file path", app.configFilePath)
		} else {
			log.Println("unable to automatically identify any suitable config dirs; configuration will not be saved")
		}
	}
}

func (app *App) saveConfig() error {
	if app.configFilePath == "" {
		return fmt.Errorf("received empty config filename")
	}

	if app.conf == nil {
		return fmt.Errorf("config was nil")
	}

	b, err := json.Marshal(app.conf)
	if err != nil {
		return fmt.Errorf("failed to marshal app config to json: %v", err.Error())
	} else {
		dir, _ := filepath.Split(app.configFilePath)
		err := os.MkdirAll(dir, 0o755)
		if err != nil {
			return fmt.Errorf("failed to create app config parent dir %v: %v", dir, err.Error())
		}
		err = os.WriteFile(app.configFilePath, b, 0o644)
		if err != nil {
			return fmt.Errorf("failed to save app config to %v: %v", app.configFilePath, err.Error())
		}
	}

	return nil
}

func (app *App) initializeInvestments() {
	if app.inv.Accounts == nil {
		app.inv.Accounts = make(map[string]balancer.Account)
	}

	if len(app.inv.Accounts) == 0 {
		app.inv.Accounts[INVESTMENTS_NAME] = balancer.Account{
			Name:     INVESTMENTS_NAME,
			Balance:  decimal.NewFromInt(5000),
			Strategy: INVESTMENTS_NAME,
		}
	}

	if app.inv.Strategies == nil {
		app.inv.Strategies = make(map[string]balancer.Strategy)
	}

	if len(app.inv.Strategies) == 0 {
		app.inv.Strategies[INVESTMENTS_NAME] = balancer.Strategy{
			Name:        INVESTMENTS_NAME,
			Symbols:     []balancer.Symbol{},
			Allocations: map[string]decimal.Decimal{},
		}
	}
}

func (app *App) repopulateSymbolsInput() {
	app.ui.symbolsInput.Clear()
	for _, symbol := range app.inv.Strategies[INVESTMENTS_NAME].Symbols {
		app.ui.symbolsInputDropdown.Add(symbol.Symbol, func() {
			app.ui.symbolsInput.SetValue(symbol.Symbol)
			app.ui.symbolEdit.SetValue(symbol.Type)
			app.ui.symbolEdit.TakeFocus()
		})
	}
}

func (app *App) repopulateAllocationsInput() {
	app.ui.allocationsInput.Clear()
	for allocation, amount := range app.inv.Strategies[INVESTMENTS_NAME].Allocations {
		app.ui.allocationsInputDropdown.Add(allocation, func() {
			app.ui.allocationsInput.SetValue(allocation)
			app.ui.allocationEdit.SetValue(amount.StringFixed(2))
			app.ui.allocationEdit.TakeFocus()
		})
		app.ui.symbolEditDropdown.Add(allocation, func() {
			app.ui.symbolEdit.SetValue(allocation)
			app.symbolEditChanged()
			app.ui.allocationsInput.SetValue(allocation)
			app.ui.allocationEdit.SetValue(amount.StringFixed(2))
			app.allocationEditChanged()
		})
	}
}

func (app *App) deleteSymbol(symbol string) {
	symbolToDelete := -1
	for i, s := range app.inv.Strategies[INVESTMENTS_NAME].Symbols {
		if s.Symbol == symbol {
			symbolToDelete = int(i)
			break
		}
	}

	if symbolToDelete == -1 {
		return
	}

	strat := app.inv.Strategies[INVESTMENTS_NAME]
	strat.Symbols = append(
		strat.Symbols[:symbolToDelete],
		strat.Symbols[symbolToDelete+1:]...,
	)

	app.inv.Strategies[INVESTMENTS_NAME] = strat

	app.ui.symbolsInput.SetValue("")
	app.ui.symbolEdit.SetValue("")
	app.repopulateSymbolsInput()
}

func (app *App) deleteAllocation(allocation string) {
	delete(app.inv.Strategies[INVESTMENTS_NAME].Allocations, allocation)
	app.ui.allocationsInput.SetValue("")
	app.ui.allocationEdit.SetValue("")
	app.repopulateAllocationsInput()
}
