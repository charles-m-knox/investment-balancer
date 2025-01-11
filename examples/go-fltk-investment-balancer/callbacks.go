package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/charles-m-knox/investment-balancer/pkg/balancer"
	"github.com/pwiecz/go-fltk"
	"github.com/shopspring/decimal"
)

func (app *App) setCallbacks() {
	app.ui.win.SetCallback(app.gracefulExit)
	app.ui.win.SetResizeHandler(func() { app.ui.responsive() })

	app.ui.menu.AddEx("Save", fltk.CTRL+'s', app.saveConfigShortcut, 0)
	app.ui.menu.AddEx("Quit", fltk.CTRL+'q', app.gracefulExit, 0)
	app.ui.menu.AddEx("Help", fltk.F1, app.help, 0)

	app.setDarkMode()
	app.setMainPage()
	app.setResultsPage()
	app.setSubmit()
	app.setSymbolsEditChanged()
	app.setSymbolEditChanged()
	app.setAllocationEditChanged()
	app.setDeleteSymbol()
	app.setDeleteAllocation()
	app.setApiKeyChanged()
	app.setBalanceChanged()
}

func (app *App) help() {
	fltk.MessageBox("Help", `Specify a ticker symbol an assign it a group (called an allocation), such as "large" or "small".

Then choose what percentage of your total balance will be assigned to each group, such as 25% for the "large" group, 10% for "small", etc.

Provide the total balance for your account and an AlphaVantage API key in the Results page and press submit.

Quote prices are cached for a while, so don't be afraid to adjust repeatedly.

Keyboard shortcuts:
Ctrl+Q: Quit
F1: Help`)
}

func (app *App) saveConfigShortcut() {
	err := app.saveConfig()
	if err != nil {
		msg := fmt.Sprintf("Failed to save: %v", err.Error())
		log.Println(msg)
		fltk.MessageBox("Error Saving", msg)
	} else {
		fltk.MessageBox("Saved", "Successfully saved configuration.")
	}
}

func (app *App) setDarkMode() {
	app.ui.dark.SetCallback(func() {
		app.conf.DarkMode = !app.conf.DarkMode
		if !app.ui.darkModeChanged {
			fltk.MessageBox(
				"App Restart Required",
				"In order for this to take effect, the theme change won't take effect until the application is restarted.",
			)

			app.ui.darkModeChanged = true
		}
		app.ui.dark.SetValue(app.conf.DarkMode)
	})
}

func (app *App) setDeleteSymbol() {
	app.ui.deleteSymbol.SetCallback(func() {
		app.deleteSymbol(app.ui.symbolsInput.Value())
		app.repopulateSymbolsInput()
	})
}

func (app *App) setDeleteAllocation() {
	app.ui.deleteAllocation.SetCallback(func() {
		app.deleteAllocation(app.ui.allocationsInput.Value())
		app.repopulateAllocationsInput()
	})
}

func (app *App) setMainPage() {
	app.ui.backBtn.SetCallback(func() {
		app.ui.page = PAGE_MAIN
		app.ui.viewPage()
	})
}

func (app *App) setResultsPage() {
	app.ui.results.SetCallback(func() {
		app.ui.page = PAGE_RESULTS
		app.ui.viewPage()
	})
}

func (app *App) setSubmit() {
	app.ui.submit.SetCallback(func() {
		err := app.saveConfig()
		if err != nil {
			log.Printf("failed to save config before submit: %v", err.Error())
		}

		quotes := app.inv.GetQuotes()
		out, err := app.inv.PrettyPrint(quotes)
		if err != nil {
			fltk.MessageBox("Error", fmt.Sprintf("Failed to process quotes: %v", err.Error()))
		}

		app.ui.output.SetBuffer(fltk.NewTextBuffer())
		app.ui.output.InsertText(out)
	})
}

func (app *App) setApiKeyChanged() {
	app.ui.apiKey.SetCallback(func() {
		app.inv.AlphaVantageAPIKey = app.ui.apiKey.Value()
	})
}

func (app *App) setBalanceChanged() {
	app.ui.totalBalance.SetCallback(func() {
		bal, err := strconv.ParseFloat(app.ui.totalBalance.Value(), 64)
		if err != nil {
			return
		}

		acct := app.inv.Accounts[INVESTMENTS_NAME]
		acct.Balance = decimal.NewFromFloat(bal)
		app.inv.Accounts[INVESTMENTS_NAME] = acct
		app.ui.totalBalance.SetValue(app.inv.Accounts[INVESTMENTS_NAME].Balance.StringFixed(2))
	})
}

func (app *App) setSymbolsEditChanged() {
	app.ui.symbolsInput.SetCallback(func() {
		v := app.ui.symbolsInput.Value()
		if v == "" {
			return
		}

		for _, symbol := range app.inv.Strategies[INVESTMENTS_NAME].Symbols {
			if v == symbol.Symbol && symbol.Type != "" {
				app.ui.symbolEdit.SetValue(symbol.Type)
				app.ui.allocationsInput.SetValue(symbol.Type)
				allocation, ok := app.inv.Strategies[INVESTMENTS_NAME].Allocations[symbol.Type]
				if ok {
					app.ui.allocationEdit.SetValue(allocation.StringFixed(2))
				}
				break
			}
		}
	})
}

func (app *App) symbolEditChanged() {
	// attempt to retrieve the user-typed symbol from the existing list of
	// symbols. If it's not there, then it's a new one and should be added
	// to the list, and the symbols edit field should be focused
	symbolName := app.ui.symbolsInput.Value()
	if symbolName == "" {
		return
	}

	allocation := app.ui.symbolEdit.Value()

	exists := false
	for i := range app.inv.Strategies[INVESTMENTS_NAME].Symbols {
		if app.inv.Strategies[INVESTMENTS_NAME].Symbols[i].Symbol == symbolName {
			app.inv.Strategies[INVESTMENTS_NAME].Symbols[i].Type = allocation
			exists = true
			break
		}
	}

	if !exists {
		strat := app.inv.Strategies[INVESTMENTS_NAME]
		strat.Symbols = append(app.inv.Strategies[INVESTMENTS_NAME].Symbols, balancer.Symbol{
			Symbol: symbolName,
			Type:   allocation,
		})
		app.inv.Strategies[INVESTMENTS_NAME] = strat
	}

	existing, ok := app.inv.Strategies[INVESTMENTS_NAME].Allocations[allocation]
	if ok {
		app.ui.allocationsInput.SetValue(allocation)
		app.ui.allocationEdit.SetValue(existing.StringFixed(2))
		defer app.repopulateAllocationsInput()
	}

	app.repopulateSymbolsInput()
}

func (app *App) allocationEditChanged() {
	// attempt to retrieve the user-typed allocation from the existing list of
	// allocations. If it's not there, then it's a new one and should be added
	// to the list, and the allocations edit field should be focused
	alloc := app.ui.allocationsInput.Value()
	if alloc == "" {
		return
	}

	v, err := strconv.ParseFloat(app.ui.allocationEdit.Value(), 64)
	if err != nil {
		log.Printf("failed to parse allocation amount int: %v", err.Error())
		return
	}

	newVal := decimal.NewFromFloat(v)

	app.inv.Strategies[INVESTMENTS_NAME].Allocations[alloc] = newVal

	app.ui.allocationEdit.SetValue(newVal.StringFixed(2))

	app.repopulateAllocationsInput()
}

func (app *App) setSymbolEditChanged() {
	app.ui.symbolEdit.SetCallback(func() {
		app.symbolEditChanged()
	})
}

func (app *App) setAllocationEditChanged() {
	app.ui.allocationEdit.SetCallback(func() {
		app.allocationEditChanged()
	})
}

func (app *App) gracefulExit() {
	log.Println("closing app and saving config, please wait a moment...")

	err := app.inv.SaveConfig()
	if err != nil {
		log.Printf("failed to save investment config: %v", err.Error())
	}

	err = app.saveConfig()
	if err != nil {
		log.Printf("failed to save config: %v", err.Error())
	}

	log.Println("done, exiting now.")
	os.Exit(0)
}
