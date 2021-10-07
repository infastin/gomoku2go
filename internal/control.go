package gomoku

import (
	"os"

	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

const (
	appID = "com.github.infastin.gomoku2go"
)

type Application struct {
	*gtk.Application
}

func NewApplication() *Application {
	app := &Application{}
	app.Application = gtk.NewApplication(appID, gio.ApplicationFlagsNone)

	return app
}

func (app *Application) Start() {
	app.Connect("activate", app.activate)

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

func (app *Application) activate(g *gtk.Application) {
	window := NewMainWindow(app)
	window.Show()
}
