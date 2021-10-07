package gomoku

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type MainWindow struct {
	*gtk.ApplicationWindow
}

func NewMainWindow(app *Application) *MainWindow {
	mwin := &MainWindow{}
	mwin.ApplicationWindow = gtk.NewApplicationWindow(app.Application)

	mwin.SetTitle("PogChamp")
	mwin.SetDefaultSize(200, 200)

	return mwin
}
