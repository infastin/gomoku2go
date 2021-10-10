package gomoku

import (
	"fmt"
	"os"

	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"

	"github.com/infastin/gomoku2go/internal/gomoku/game"
	"github.com/infastin/gomoku2go/internal/gomoku/view"
)

const (
	appID = "com.github.infastin.gomoku2go"
)

type Application struct {
	*gtk.Application

	settings *gio.Settings

	gameView  *view.MainWindow
	gameLogic *game.Game
}

func NewApplication() *Application {
	app := &Application{}
	app.Application = gtk.NewApplication(appID, gio.ApplicationFlagsNone)
	return app
}

func (app *Application) handleClick(x, y uint) {
	suc, _ := app.gameLogic.SetField(x, y)

	board := app.gameView.Board()

	if suc {
		switch app.gameLogic.CurrentPlayer() {
		case game.FirstPlayer:
			glib.IdleAdd(func() {
				board.DrawCircle(x, y)
			})
		case game.SecondPlayer:
			glib.IdleAdd(func() {
				board.DrawCross(x, y)
			})
		}

		app.gameLogic.ChangePlayer()
	}
}

func (app *Application) handleRedraw() {
	lfields := app.gameLogic.NotEmptyFields()
	board := app.gameView.Board()

	var vfields []view.Field

	for _, lf := range lfields {
		var sh view.Shape

		switch lf.Ft {
		case game.FirstPlayerField:
			sh = view.Circle
		case game.SecondPlayerField:
			sh = view.Cross
		}

		vfields = append(vfields, view.Field{
			X:  lf.X,
			Y:  lf.Y,
			Sh: sh,
		})
	}

	glib.IdleAdd(func() {
		board.DrawShapes(vfields)
	})
}

func (app *Application) Start() {
	app.ConnectActivate(app.activate)
	app.ConnectStartup(app.startup)

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

func (app *Application) activate() {
	app.gameView = view.NewMainWindow(app.Application)
	app.gameView.Show()

	board := app.gameView.Board()

	board.DrawBoard(app.gameLogic.Size())
	board.ConnectClick(app.handleClick)
	board.ConnectRedraw(app.handleRedraw)
}

func (app *Application) quit() {
	app.Quit()
}

func (app *Application) prefs() {
	fmt.Println("To be done")
}

func (app *Application) startup() {
	app.AddAction(NewAction("preferences", nil, app.prefs))
	app.AddAction(NewAction("quit", nil, app.quit))

	app.settings = gio.NewSettings(appID)

	p1Str := app.settings.String("player1")
	p2Str := app.settings.String("player2")

	if len(p1Str) > 16 {
		p1Str = "Player 1"
		app.settings.SetString("player1", p1Str)
	}

	if len(p2Str) > 16 {
		p2Str = "Player 2"
		app.settings.SetString("player2", p2Str)
	}

	size := app.settings.Uint("size")
	wincond := app.settings.Uint("wincond")

	p1 := game.NewPlayer(p1Str)
	p2 := game.NewPlayer(p2Str)

	app.gameLogic, _ = game.NewGame(p1, p2, size, wincond)
}
