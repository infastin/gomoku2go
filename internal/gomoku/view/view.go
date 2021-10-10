package view

import (
	_ "embed"

	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

//go:embed resources/window.ui
var winui string

//go:embed resources/titlebar.ui
var titleui string

//go:embed resources/menu.ui
var menui string

type MainWindow struct {
	*gtk.ApplicationWindow

	board *BoardArea
}

type NeedRedraw struct{}

func NewMainWindow(app *gtk.Application) *MainWindow {
	mwin := &MainWindow{}

	builder := gtk.NewBuilderFromString(winui, len(winui))
	builderBar := gtk.NewBuilderFromString(titleui, len(titleui))
	builderMenu := gtk.NewBuilderFromString(menui, len(menui))

	mwin.ApplicationWindow = builder.GetObject("mainwin").Cast().(*gtk.ApplicationWindow)

	titlebar := builderBar.GetObject("titlebar").Cast().(*gtk.HeaderBar)
	mwin.SetTitlebar(titlebar)

	menuButton := builderBar.GetObject("menu_button").Cast().(*gtk.MenuButton)
	menu := builderMenu.GetObject("menu").Cast().(*gio.Menu)
	menuButton.SetMenuModel(menu)

	mwin.board = newBoardArea(builder)

	mwin.SetApplication(app)

	return mwin
}

func (mwin *MainWindow) Board() *BoardArea {
	return mwin.board
}
