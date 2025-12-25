package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

type App struct {
	window   fyne.Window
	parserUI *ParserUI
}

func NewApp() *App {
	myApp := app.New()
	myApp.Settings().SetTheme(theme.DarkTheme())
	
	app := &App{
		window: myApp.NewWindow("Pars Channel by @MaxGolang"),
	}

	app.parserUI = NewParserUI()
	app.window.SetContent(app.parserUI.BuildUI())

	if w, h := loadWindowSize(); w >= 600 && h >= 400 {
		app.window.Resize(fyne.NewSize(float32(w), float32(h)))
	} else {
		app.window.Resize(fyne.NewSize(900, 520))
	}
	app.window.CenterOnScreen()

	app.window.SetCloseIntercept(func() {
		saveWindowSize(app.window.Canvas().Size())
		app.window.Close()
	})
	
	return app
}

func (a *App) Run() {
	a.window.ShowAndRun()
}

