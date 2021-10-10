package view

import (
	_ "embed"
	"fmt"

	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotk4/pkg/pango"
)

//go:embed resources/window.ui
var winui string

//go:embed resources/titlebar.ui
var titleui string

//go:embed resources/menu.ui
var menui string

//go:embed resources/style.css
var style string

type MainWindow struct {
	*gtk.ApplicationWindow

	board          *BoardArea
	startGame      *gtk.Button
	curPlayerLabel *gtk.Label
	stopwatch      *Stopwatch
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

	mwin.startGame = builderBar.GetObject("start_game").Cast().(*gtk.Button)

	menuButton := builderBar.GetObject("menu_button").Cast().(*gtk.MenuButton)
	menu := builderMenu.GetObject("menu").Cast().(*gio.Menu)
	menuButton.SetMenuModel(menu)

	mwin.curPlayerLabel = builder.GetObject("current_player").Cast().(*gtk.Label)
	stopwatch := builder.GetObject("stopwatch").Cast().(*gtk.Label)
	mwin.stopwatch = NewStopwatch(stopwatch)

	mwin.board = newBoardArea(builder)

	css := gtk.NewCSSProvider()
	css.LoadFromData(style)

	scxt := mwin.StyleContext()
	gtk.StyleContextAddProviderForDisplay(scxt.Display(), css, gtk.STYLE_PROVIDER_PRIORITY_USER)

	fmt.Println(mwin.PangoContext().FontDescription().Size(), pango.SCALE)

	mwin.SetApplication(app)

	return mwin
}

func (mwin *MainWindow) Board() *BoardArea {
	return mwin.board
}

func (mwin *MainWindow) CurrentPlayerLabel() *gtk.Label {
	return mwin.curPlayerLabel
}

func (mwin *MainWindow) Stopwatch() *Stopwatch {
	return mwin.stopwatch
}

func (mwin *MainWindow) StartGameBtn() *gtk.Button {
	return mwin.startGame
}
