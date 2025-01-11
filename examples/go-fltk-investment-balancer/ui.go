package main

import (
	"fmt"
	"log"
	"math"

	"github.com/pwiecz/go-fltk"
)

// If the screen is portrait or landscape, the window will be scaled
// accordingly.
const (
	WIDTH_PORTRAIT   = 100
	HEIGHT_PORTRAIT  = 150
	WIDTH_LANDSCAPE  = 150
	HEIGHT_LANDSCAPE = 100
)

const (
	PAGE_MAIN = iota
	PAGE_RESULTS
)

// Positioning (x,y,w,h) for fltk elements
type pos struct {
	X int
	Y int
	W int
	H int
	// A reference to the UI is needed in order to perform the translation
	// according to the ui's specs
	ui *UI
}

// Buttons, inputs, widgets, etc that need to be repositioned in a responsive
// manner.
type UI struct {
	win *fltk.Window // main window

	menu *fltk.MenuBar // hidden menu bar for shortcut keys

	page int

	// main page components & positions
	symbolsInput             *fltk.InputChoice // symbol list
	symbolsInputDropdown     *fltk.MenuButton  // symbol list
	symbolsInputp            pos               // scrollable symbol list
	allocationsInput         *fltk.InputChoice // allocations list
	allocationsInputDropdown *fltk.MenuButton  // allocations list
	allocationsInputp        pos               // allocations list
	symbolEdit               *fltk.InputChoice // edits the group for a symbol
	symbolEditDropdown       *fltk.MenuButton  // shows all the user's allocations to help with ux
	symbolEditp              pos               // symbol edit position
	allocationEdit           *fltk.Input       // edits the percentage for an allocation
	allocationEditp          pos               // allocation edit position
	deleteSymbol             *fltk.Button      // allocations page button
	deleteSymbolp            pos               // allocations page button position
	deleteAllocation         *fltk.Button      // allocation deletion button
	deleteAllocationp        pos               // allocation deletion button position
	totalBalance             *fltk.Input       // stores the total balance for the account
	totalBalancep            pos               // position of the total balance input
	results                  *fltk.Button      // results page button
	resultsp                 pos               // results page button position

	// results page components & positions
	dark     *fltk.CheckButton // dark mode toggle
	darkp    pos               // dark mode toggle position
	output   *fltk.TextDisplay // text field that displays results
	outputp  pos               // text field that displays results
	submit   *fltk.Button      // submit button for results page
	submitp  pos               // submit button for results page position
	apiKey   *fltk.Input       // api key input field
	apiKeyp  pos               // api key input field position
	backBtn  *fltk.Button      // symbols page button
	backBtnp pos               // symbols page button position

	portrait        bool // portrait mode or landscape mode
	darkModeChanged bool // if true, prompts to restart after changing dark mode will not show
}

// isPortrait returns true if the screen is taller than it is wide. It returns
// false otherwise, including for square screens.
func isPortrait() (bool, error) {
	_, _, width, height := fltk.ScreenWorkArea(int(fltk.SCREEN))

	if width == 0 || height == 0 {
		return false, fmt.Errorf("received 0 for one of screen height or width")
	}

	if width > height {
		return false, nil
	}

	return true, nil
}

// Translates the widget's width/height from the original 100 or 150px base
// width/height to the window's current width/height
func (ui *UI) tr(i int, winw int, winh int, useHeight bool) int {
	if ui.portrait {
		if useHeight {
			return int(math.Round((float64(i) / float64(HEIGHT_PORTRAIT)) * float64(winh)))
		} else {
			return int(math.Round((float64(i) / float64(WIDTH_PORTRAIT)) * float64(winw)))
		}
	} else {
		if useHeight {
			return int(math.Round((float64(i) / float64(HEIGHT_LANDSCAPE)) * float64(winh)))
		} else {
			return int(math.Round((float64(i) / float64(WIDTH_LANDSCAPE)) * float64(winw)))
		}
	}
}

// Translate converts a predefined position into a scaled position based on
// the latest width & height of the window.
func (p *pos) Translate(winw, winh int) {
	p.X = p.ui.tr(p.X, winw, winh, false)
	p.Y = p.ui.tr(p.Y, winw, winh, true)
	p.W = p.ui.tr(p.W, winw, winh, false)
	p.H = p.ui.tr(p.H, winw, winh, true)
}

// Initializes the UI for the app. Call this once, only after the app config has
// been loaded.
func (app *App) initUI() {
	// fltk.SetScheme("gtk+")
	fltk.InitStyles()
	fltk.SetTooltipDelay(0.1)
	fltk.EnableTooltips()

	var err error
	app.ui.portrait, err = isPortrait()
	if err != nil {
		log.Fatalf("failed to determine screen size: %v", err.Error())
	}

	// probably could write this more intelligently later
	winw := WIDTH_LANDSCAPE
	winh := HEIGHT_LANDSCAPE
	if app.ui.portrait || forcePortrait {
		winw = WIDTH_PORTRAIT
		winh = HEIGHT_PORTRAIT
		app.ui.portrait = true
	}

	if forceLandscape {
		winw = WIDTH_LANDSCAPE
		winh = HEIGHT_LANDSCAPE
		app.ui.portrait = false
	}

	// initialize all buttons and widgets
	app.ui.win = fltk.NewWindow(winw, winh, "Investment Balancer FLTK")
	app.ui.menu = fltk.NewMenuBar(0, 0, 0, 0)

	// main page components
	app.ui.symbolsInput = fltk.NewInputChoice(0, 0, 0, 0, "Symbol")
	app.ui.symbolEdit = fltk.NewInputChoice(0, 0, 0, 0, "Allocation (small/broad/etc)")
	app.ui.allocationsInput = fltk.NewInputChoice(0, 0, 0, 0, "Allocations")
	app.ui.allocationEdit = fltk.NewInput(0, 0, 0, 0, "Specify percent of total")
	app.ui.deleteSymbol = fltk.NewButton(0, 0, 0, 0, "Delete Symbol")
	app.ui.deleteAllocation = fltk.NewButton(0, 0, 0, 0, "Delete Allocation")
	app.ui.totalBalance = fltk.NewInput(0, 0, 0, 0, "&Total Balance")
	app.ui.results = fltk.NewButton(0, 0, 0, 0, "&Results")

	// results page components
	app.ui.dark = fltk.NewCheckButton(0, 0, 0, 0, "&Dark Mode")
	app.ui.output = fltk.NewTextDisplay(0, 0, 0, 0, "")
	app.ui.submit = fltk.NewButton(0, 0, 0, 0, "&Submit")
	app.ui.apiKey = fltk.NewInput(0, 0, 0, 0, "API &Key")
	app.ui.backBtn = fltk.NewButton(0, 0, 0, 0, "&Back")

	// propagate default values from config to widgets that accept them
	app.ui.dark.SetValue(app.conf.DarkMode)
	app.ui.totalBalance.SetValue(app.inv.Accounts[INVESTMENTS_NAME].Balance.StringFixed(2))
	app.ui.apiKey.SetValue(app.inv.AlphaVantageAPIKey)

	// app.ui.symbolEdit.SetCallbackCondition(fltk.WhenEnterKey)
	// app.ui.symbolsInput.SetCallbackCondition(fltk.WhenRelease)
	// app.ui.symbolsInput.SetCallback(func() { log.Println("changed", app.ui.symbolsInput.Value()) })
	// app.ui.allocationsInput.SetCallbackCondition(fltk.WhenRelease)
	// app.ui.allocationsInput.SetCallback(func() { log.Println("changed", app.ui.allocationsInput.Value()) })

	app.ui.symbolsInputDropdown = app.ui.symbolsInput.MenuButton()
	app.repopulateSymbolsInput()
	app.ui.allocationsInputDropdown = app.ui.allocationsInput.MenuButton()
	app.ui.symbolEditDropdown = app.ui.symbolEdit.MenuButton()
	app.repopulateAllocationsInput()

	app.ui.symbolsInput.SetAlign(fltk.ALIGN_TOP_LEFT)
	app.ui.symbolEdit.SetAlign(fltk.ALIGN_TOP_LEFT)
	app.ui.allocationsInput.SetAlign(fltk.ALIGN_TOP_LEFT)
	app.ui.allocationEdit.SetAlign(fltk.ALIGN_TOP_LEFT)
	app.ui.apiKey.SetAlign(fltk.ALIGN_TOP_LEFT)

	app.ui.output.SetTextFont(fltk.COURIER)

	app.ui.allocationEdit.SetType(1) // FL_FLOAT_INPUT
	app.ui.totalBalance.SetType(1)   // FL_FLOAT_INPUT
	app.ui.apiKey.SetType(5)         // FL_SECRET_INPUT

	app.ui.win.Resizable(app.ui.win)
	app.ui.win.SetXClass("gfltkinvbal")
}

// Sizes the window to 3x the design size, which is intentionally
// small. Use this after things have been initiated.
func (ui *UI) upsize() {
	if app.ui.portrait {
		app.ui.win.Resize(0, 0, WIDTH_PORTRAIT*3, HEIGHT_PORTRAIT*3)
	} else {
		app.ui.win.Resize(0, 0, WIDTH_LANDSCAPE*3, HEIGHT_LANDSCAPE*3)
	}
}

// Resizes and repositions all components based on the window's size.
func (ui *UI) responsive() {
	if forceLandscape || forcePortrait {
		return
	}

	winw := ui.win.W()
	winh := ui.win.H()

	if winw > winh {
		ui.portrait = false
	} else {
		ui.portrait = true
	}

	if ui.portrait {
		// main page
		ui.symbolsInputp = pos{X: 5, Y: 10, W: 90, H: 10, ui: ui}
		ui.symbolEditp = pos{X: 5, Y: 30, W: 90, H: 10, ui: ui}
		ui.allocationsInputp = pos{X: 5, Y: 70, W: 90, H: 10, ui: ui}
		ui.allocationEditp = pos{X: 5, Y: 90, W: 90, H: 10, ui: ui}
		ui.deleteSymbolp = pos{X: 5, Y: 45, W: 90, H: 10, ui: ui}
		ui.deleteAllocationp = pos{X: 5, Y: 105, W: 90, H: 10, ui: ui}
		ui.totalBalancep = pos{X: 35, Y: 120, W: 60, H: 10, ui: ui}
		ui.resultsp = pos{X: 5, Y: 135, W: 90, H: 10, ui: ui}

		// results page
		ui.outputp = pos{X: 5, Y: 5, W: 90, H: 95, ui: ui}
		ui.darkp = pos{X: 60, Y: 105, W: 35, H: 10, ui: ui}
		ui.apiKeyp = pos{X: 5, Y: 120, W: 60, H: 10, ui: ui}
		ui.submitp = pos{X: 70, Y: 120, W: 25, H: 10, ui: ui}
		ui.backBtnp = pos{X: 5, Y: 135, W: 90, H: 10, ui: ui}
	} else {
		// landscape

		// main page
		ui.symbolsInputp = pos{X: 5, Y: 10, W: 65, H: 10, ui: ui}
		ui.symbolEditp = pos{X: 5, Y: 30, W: 65, H: 10, ui: ui}
		ui.allocationsInputp = pos{X: 80, Y: 10, W: 65, H: 10, ui: ui}
		ui.allocationEditp = pos{X: 80, Y: 30, W: 65, H: 10, ui: ui}
		ui.deleteSymbolp = pos{X: 5, Y: 50, W: 65, H: 10, ui: ui}
		ui.deleteAllocationp = pos{X: 80, Y: 50, W: 65, H: 10, ui: ui}
		ui.totalBalancep = pos{X: 45, Y: 70, W: 60, H: 10, ui: ui}
		ui.resultsp = pos{X: 5, Y: 85, W: 140, H: 10, ui: ui}

		// results page
		ui.outputp = pos{X: 5, Y: 5, W: 140, H: 45, ui: ui}
		ui.darkp = pos{X: 110, Y: 55, W: 35, H: 10, ui: ui}
		ui.apiKeyp = pos{X: 5, Y: 70, W: 110, H: 10, ui: ui}
		ui.submitp = pos{X: 120, Y: 70, W: 25, H: 10, ui: ui}
		ui.backBtnp = pos{X: 5, Y: 85, W: 140, H: 10, ui: ui}
	}
	// main page translation of coordinates to newly resized window
	ui.symbolsInputp.Translate(winw, winh)
	ui.symbolEditp.Translate(winw, winh)
	ui.allocationsInputp.Translate(winw, winh)
	ui.allocationEditp.Translate(winw, winh)
	ui.deleteSymbolp.Translate(winw, winh)
	ui.deleteAllocationp.Translate(winw, winh)
	ui.totalBalancep.Translate(winw, winh)
	ui.resultsp.Translate(winw, winh)
	// results page translation of coordinates to newly resized window
	ui.outputp.Translate(winw, winh)
	ui.darkp.Translate(winw, winh)
	ui.apiKeyp.Translate(winw, winh)
	ui.submitp.Translate(winw, winh)
	ui.backBtnp.Translate(winw, winh)
	// main page resizing
	ui.symbolsInput.Resize(ui.symbolsInputp.X, ui.symbolsInputp.Y, ui.symbolsInputp.W, ui.symbolsInputp.H)
	ui.symbolEdit.Resize(ui.symbolEditp.X, ui.symbolEditp.Y, ui.symbolEditp.W, ui.symbolEditp.H)
	ui.allocationsInput.Resize(ui.allocationsInputp.X, ui.allocationsInputp.Y, ui.allocationsInputp.W, ui.allocationsInputp.H)
	ui.allocationEdit.Resize(ui.allocationEditp.X, ui.allocationEditp.Y, ui.allocationEditp.W, ui.allocationEditp.H)
	ui.deleteSymbol.Resize(ui.deleteSymbolp.X, ui.deleteSymbolp.Y, ui.deleteSymbolp.W, ui.deleteSymbolp.H)
	ui.deleteAllocation.Resize(ui.deleteAllocationp.X, ui.deleteAllocationp.Y, ui.deleteAllocationp.W, ui.deleteAllocationp.H)
	ui.totalBalance.Resize(ui.totalBalancep.X, ui.totalBalancep.Y, ui.totalBalancep.W, ui.totalBalancep.H)
	ui.results.Resize(ui.resultsp.X, ui.resultsp.Y, ui.resultsp.W, ui.resultsp.H)
	// results page resizing
	ui.output.Resize(ui.outputp.X, ui.outputp.Y, ui.outputp.W, ui.outputp.H)
	ui.dark.Resize(ui.darkp.X, ui.darkp.Y, ui.darkp.W, ui.darkp.H)
	ui.apiKey.Resize(ui.apiKeyp.X, ui.apiKeyp.Y, ui.apiKeyp.W, ui.apiKeyp.H)
	ui.submit.Resize(ui.submitp.X, ui.submitp.Y, ui.submitp.W, ui.submitp.H)
	ui.backBtn.Resize(ui.backBtnp.X, ui.backBtnp.Y, ui.backBtnp.W, ui.backBtnp.H)
}

const (
	DARK_COLOR_TEXT               fltk.Color = 0x9f9f9f00
	DARK_COLOR_INPUT_BG           fltk.Color = 0x20202000
	DARK_COLOR_INPUT_SELECTED_BG  fltk.Color = 0xafafaf00
	LIGHT_COLOR_TEXT              fltk.Color = 0x20030500
	LIGHT_COLOR_INPUT_BG          fltk.Color = 0xFFFFFF00
	LIGHT_COLOR_INPUT_SELECTED_BG fltk.Color = 0x00008000
)

var (
	COLOR_TEXT              fltk.Color = LIGHT_COLOR_TEXT
	COLOR_INPUT_BG          fltk.Color = LIGHT_COLOR_INPUT_BG
	COLOR_INPUT_SELECTED_BG fltk.Color = LIGHT_COLOR_INPUT_SELECTED_BG
)

// Changes the color of various widgets/states to light or dark mode.
func (ui *UI) theme(dark bool) {
	if dark {
		log.Println("dark mode activated")
		COLOR_TEXT = DARK_COLOR_TEXT
		COLOR_INPUT_BG = DARK_COLOR_INPUT_BG
		COLOR_INPUT_SELECTED_BG = DARK_COLOR_INPUT_SELECTED_BG
		fltk.SetForegroundColor(230, 230, 230)
		fltk.SetBackgroundColor(40, 40, 40)
	} else {
		log.Println("light mode activated")
		COLOR_TEXT = LIGHT_COLOR_TEXT
		COLOR_INPUT_BG = LIGHT_COLOR_INPUT_BG
		COLOR_INPUT_SELECTED_BG = LIGHT_COLOR_INPUT_SELECTED_BG
		fltk.SetBackgroundColor(192, 192, 192)
		fltk.SetForegroundColor(0, 0, 0)
		return
	}

	// main page text color
	ui.allocationsInput.SetLabelColor(COLOR_TEXT)
	ui.allocationsInput.Input().SetLabelColor(COLOR_TEXT)
	ui.allocationEdit.SetLabelColor(COLOR_TEXT)
	ui.symbolsInput.SetLabelColor(COLOR_TEXT)
	ui.symbolsInput.Input().SetLabelColor(COLOR_TEXT)
	ui.symbolEdit.SetLabelColor(COLOR_TEXT)
	ui.symbolEdit.Input().SetLabelColor(COLOR_TEXT)
	ui.deleteSymbol.SetLabelColor(COLOR_TEXT)
	ui.deleteAllocation.SetLabelColor(COLOR_TEXT)
	ui.totalBalance.SetLabelColor(COLOR_TEXT)
	ui.results.SetLabelColor(COLOR_TEXT)
	// results page text color
	ui.output.SetLabelColor(COLOR_TEXT)
	ui.dark.SetLabelColor(COLOR_TEXT)
	ui.apiKey.SetLabelColor(COLOR_TEXT)
	ui.submit.SetLabelColor(COLOR_TEXT)
	ui.backBtn.SetLabelColor(COLOR_TEXT)
	// main page input color
	ui.allocationsInput.SetColor(COLOR_INPUT_BG)
	ui.allocationsInput.Input().SetColor(COLOR_INPUT_BG)
	ui.allocationEdit.SetColor(COLOR_INPUT_BG)
	ui.symbolsInput.SetColor(COLOR_INPUT_BG)
	ui.symbolsInput.Input().SetColor(COLOR_INPUT_BG)
	ui.symbolEdit.SetColor(COLOR_INPUT_BG)
	ui.symbolEdit.Input().SetColor(COLOR_INPUT_BG)
	ui.deleteSymbol.SetColor(COLOR_INPUT_BG)
	ui.deleteAllocation.SetColor(COLOR_INPUT_BG)
	ui.totalBalance.SetColor(COLOR_INPUT_BG)
	ui.results.SetColor(COLOR_INPUT_BG)
	// results page input color
	ui.output.SetColor(COLOR_INPUT_BG)
	ui.dark.SetColor(COLOR_INPUT_BG)
	ui.submit.SetColor(COLOR_INPUT_BG)
	ui.apiKey.SetColor(COLOR_INPUT_BG)
	ui.backBtn.SetColor(COLOR_INPUT_BG)
	// main page input selection color
	ui.allocationsInput.SetSelectionColor(COLOR_INPUT_SELECTED_BG)
	ui.allocationsInput.Input().SetSelectionColor(COLOR_INPUT_SELECTED_BG)
	ui.allocationEdit.SetSelectionColor(COLOR_INPUT_SELECTED_BG)
	ui.symbolsInput.SetSelectionColor(COLOR_INPUT_SELECTED_BG)
	ui.symbolsInput.Input().SetSelectionColor(COLOR_INPUT_SELECTED_BG)
	ui.symbolEdit.SetSelectionColor(COLOR_INPUT_SELECTED_BG)
	ui.symbolEdit.Input().SetSelectionColor(COLOR_INPUT_SELECTED_BG)
	ui.deleteAllocation.SetSelectionColor(COLOR_INPUT_SELECTED_BG)
	ui.deleteSymbol.SetSelectionColor(COLOR_INPUT_SELECTED_BG)
	ui.totalBalance.SetSelectionColor(COLOR_INPUT_SELECTED_BG)
	ui.results.SetSelectionColor(COLOR_INPUT_SELECTED_BG)
	// results page input selection color
	ui.output.SetSelectionColor(COLOR_INPUT_SELECTED_BG)
	ui.dark.SetSelectionColor(COLOR_INPUT_SELECTED_BG)
	ui.apiKey.SetSelectionColor(COLOR_INPUT_SELECTED_BG)
	ui.submit.SetSelectionColor(COLOR_INPUT_SELECTED_BG)
	ui.backBtn.SetSelectionColor(COLOR_INPUT_SELECTED_BG)
}

func (ui *UI) viewPage() {
	switch ui.page {
	case PAGE_MAIN:
		ui.allocationsInput.Activate()
		ui.allocationsInput.Show()
		ui.allocationEdit.Activate()
		ui.allocationEdit.Show()
		ui.symbolsInput.Activate()
		ui.symbolsInput.Show()
		ui.symbolEdit.Activate()
		ui.symbolEdit.Show()
		ui.deleteSymbol.Activate()
		ui.deleteSymbol.Show()
		ui.deleteAllocation.Activate()
		ui.deleteAllocation.Show()
		ui.totalBalance.Activate()
		ui.totalBalance.Show()
		ui.results.Activate()
		ui.results.Show()

		ui.output.Hide()
		ui.output.Deactivate()
		ui.dark.Hide()
		ui.dark.Deactivate()
		ui.apiKey.Hide()
		ui.apiKey.Deactivate()
		ui.submit.Hide()
		ui.submit.Deactivate()
		ui.backBtn.Hide()
		ui.backBtn.Deactivate()
	case PAGE_RESULTS:
		ui.allocationsInput.Deactivate()
		ui.allocationsInput.Hide()
		ui.allocationEdit.Deactivate()
		ui.allocationEdit.Hide()
		ui.symbolsInput.Deactivate()
		ui.symbolsInput.Hide()
		ui.symbolEdit.Deactivate()
		ui.symbolEdit.Hide()
		ui.deleteSymbol.Deactivate()
		ui.deleteSymbol.Hide()
		ui.deleteAllocation.Deactivate()
		ui.deleteAllocation.Hide()
		ui.totalBalance.Deactivate()
		ui.totalBalance.Hide()
		ui.results.Deactivate()
		ui.results.Hide()

		ui.output.Show()
		ui.output.Activate()
		ui.dark.Show()
		ui.dark.Activate()
		ui.apiKey.Show()
		ui.apiKey.Activate()
		ui.submit.Show()
		ui.submit.Activate()
		ui.backBtn.Show()
		ui.backBtn.Activate()
	}
}
