package gomoku

import (
	"fmt"
	"math"
	"os"

	"github.com/diamondburned/gotk4/pkg/gio/v2"
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
	suc, err := app.gameLogic.SetField(x, y)

	if err != nil {
		return
	}

	board := app.gameView.Board()

	if suc {
		switch app.gameLogic.CurrentPlayer() {
		case game.FirstPlayer:
			board.DrawCircle(x, y)
		case game.SecondPlayer:
			board.DrawCross(x, y)
		}

		app.gameLogic.CheckDraw()

		if s, win := app.gameLogic.CheckWinner(x, y); win {
			app.gameView.Stopwatch().Stop()
			app.gameView.StartGameBtn().SetLabel("Start Game")
			board.DrawStrike(s.X0, s.Y0, s.X1, s.Y1)
		}

		app.gameLogic.ChangePlayer()
		playerName := app.gameLogic.Player(app.gameLogic.CurrentPlayer()).Name()
		app.gameView.CurrentPlayerLabel().SetText(fmt.Sprintf("%s's turn", playerName))
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

	board.DrawShapes(vfields)
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
	app.gameView.StartGameBtn().ConnectClicked(app.startGame)
}

func (app *Application) startGame() {
	p1Str := app.settings.String("player1")
	p2Str := app.settings.String("player2")
	size := app.settings.Uint("size")
	winCond := app.settings.Uint("wincond")

	board := app.gameView.Board()
	btn := app.gameView.StartGameBtn()
	cplabel := app.gameView.CurrentPlayerLabel()

	p1 := game.NewPlayer(p1Str)
	p2 := game.NewPlayer(p2Str)

	if app.gameLogic != nil {
		app.gameView.Stopwatch().Reset()
		board.Clear()
	}

	app.gameLogic, _ = game.NewGame(p1, p2, size, winCond)

	playerName := app.gameLogic.Player(app.gameLogic.CurrentPlayer()).Name()
	cplabel.SetText(fmt.Sprintf("%s's turn", playerName))

	btn.SetLabel("Restart Game")

	board.Init(size)
	board.Draw()
	board.ConnectClick(app.handleClick)
	board.ConnectRedraw(app.handleRedraw)

	app.gameView.Stopwatch().Start()
}

func (app *Application) quit() {
	app.Quit()
}

func (app *Application) prefs() {
	dialog := view.NewPrefsDialog(app.gameView)
	dialog.Show()

	p1Entry := dialog.FirstPlayerEntry()
	p2Entry := dialog.SecondPlayerEntry()
	sizeSB := dialog.BoardSizeSpinButton()
	wincondSB := dialog.WinCondSpinButton()
	errorLabel := dialog.ErrorLabel()

	p1Entry.SetText(app.settings.String("player1"))
	p2Entry.SetText(app.settings.String("player2"))
	sizeSB.SetValue(float64(app.settings.Uint("size")))
	wincondSB.SetValue(float64(app.settings.Uint("wincond")))

	dialog.CancelButton().ConnectClicked(func() {
		dialog.Close()
	})

	dialog.ConfirmButton().ConnectClicked(func() {
		if errorLabel.Visible() {
			errorLabel.SetVisible(false)
		}

		p1Name := p1Entry.Text()
		p2Name := p2Entry.Text()
		size := sizeSB.Value()
		wincond := wincondSB.Value()

		usize := uint(math.Floor(size))
		uwincond := uint(math.Floor(wincond))

		p1 := game.NewPlayer(p1Name)
		p2 := game.NewPlayer(p2Name)

		err := game.CheckSettings(p1, p2,
			usize, uwincond)

		if err == nil {
			app.settings.SetString("player1", p1Name)
			app.settings.SetString("player2", p2Name)
			app.settings.SetUint("size", usize)
			app.settings.SetUint("wincond", uwincond)

			dialog.Close()
		} else {
			errorLabel.SetText(fmt.Sprint("Error: ", err.Error()))
			errorLabel.SetVisible(true)
		}
	})
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
}
