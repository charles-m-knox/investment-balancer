package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pwiecz/go-fltk"

	inv "github.com/charles-m-knox/investment-balancer/pkg/app"
)

// The version of the application; set at build time via:
//
//	`go build -ldflags "-X main.version=1.2.3" main.go`
//
//nolint:revive
var version string = "dev"

// Flag for showing the version and subsequently quitting.
var flagVersion bool

// Used for the config file directory and other things.
const (
	APP_NAME        = "go-fltk-investment-balancer"
	CONFIG_FILE     = "config.json"
	INV_CONFIG_FILE = "investments.json"
	// The name of the account & strategy that this application will use in the
	// config.
	INVESTMENTS_NAME = "FLTK"
)

var (
	// If true, the app will always be rendered in portrait mode
	forcePortrait bool
	// If true, the app will always be rendered in landscape mode
	forceLandscape bool
	// app contains the shared state that is required for the entire app to
	// function.
	app App = App{
		conf: &AppConfig{},
		inv:  &inv.Config{},
		ui:   &UI{},
	}
)

// App contains the shared state that is required for the entire app to
// function.
type App struct {
	// The configuration for the entire app, loaded and saved to the XDG config.
	conf *AppConfig
	inv  *inv.Config
	// All of the UI elements for this app are contained in the UI struct.
	ui *UI
	// Data is stored between runs of this application in this config file.
	configFilePath string
}

type AppConfig struct {
	// If true, the app will start in dark mode
	DarkMode bool `json:"darkMode"`
}

func parseFlags() {
	flag.BoolVar(&forcePortrait, "portrait", false, "force portrait orientation for the interface")
	flag.BoolVar(&forceLandscape, "landscape", false, "force landscape orientation for the interface")
	flag.StringVar(&app.configFilePath, "f", "", "the config file to write to, instead of the default provided by XDG config directories")
	flag.BoolVar(&flagVersion, "v", false, "print version and exit")
	flag.Parse()
}

func main() {
	parseFlags()
	if flagVersion {
		//nolint:forbidigo
		fmt.Println(version)
		os.Exit(0)
	}

	err := app.inv.LoadConfig(INV_CONFIG_FILE, APP_NAME, INV_CONFIG_FILE)
	if err != nil {
		log.Fatalf("failed to load config: %v", err.Error())
	}

	app.loadConfig()
	app.initializeInvestments()
	app.initUI()
	app.ui.theme(app.conf.DarkMode)
	app.ui.responsive()
	app.ui.upsize()
	app.setCallbacks()
	app.ui.win.End()
	app.ui.win.Show()
	app.ui.viewPage()
	go fltk.Run()

	// Channel that receives OS signals, like ctrl+c to interrupt
	var sc chan os.Signal = make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	// Block until a signal is received
	<-sc

	app.gracefulExit()
}
