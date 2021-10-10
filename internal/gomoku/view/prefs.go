package view

import (
	_ "embed"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

//go:embed resources/prefs.ui
var prefsui string

type PrefsDialog struct {
	*gtk.Dialog

	cancel    *gtk.Button
	confirm   *gtk.Button
	error     *gtk.Label
	wincond   *gtk.SpinButton
	boardsize *gtk.SpinButton
	p1        *gtk.Entry
	p2        *gtk.Entry
}

func NewPrefsDialog(mwin *MainWindow) *PrefsDialog {
	prefs := &PrefsDialog{}

	builder := gtk.NewBuilderFromString(prefsui, len(prefsui))
	prefs.Dialog = builder.GetObject("prefs").Cast().(*gtk.Dialog)

	prefs.SetTransientFor(&mwin.Window)

	prefs.cancel = builder.GetObject("cancel").Cast().(*gtk.Button)
	prefs.confirm = builder.GetObject("confirm").Cast().(*gtk.Button)
	prefs.error = builder.GetObject("error_label").Cast().(*gtk.Label)
	prefs.wincond = builder.GetObject("wincond_sb").Cast().(*gtk.SpinButton)
	prefs.boardsize = builder.GetObject("boardsize_sb").Cast().(*gtk.SpinButton)
	prefs.p1 = builder.GetObject("player1_entry").Cast().(*gtk.Entry)
	prefs.p2 = builder.GetObject("player2_entry").Cast().(*gtk.Entry)

	return prefs
}

func (p *PrefsDialog) CancelButton() *gtk.Button {
	return p.cancel
}

func (p *PrefsDialog) ConfirmButton() *gtk.Button {
	return p.confirm
}

func (p *PrefsDialog) ErrorLabel() *gtk.Label {
	return p.error
}

func (p *PrefsDialog) WinCondSpinButton() *gtk.SpinButton {
	return p.wincond
}

func (p *PrefsDialog) BoardSizeSpinButton() *gtk.SpinButton {
	return p.boardsize
}

func (p *PrefsDialog) FirstPlayerEntry() *gtk.Entry {
	return p.p1
}

func (p *PrefsDialog) SecondPlayerEntry() *gtk.Entry {
	return p.p2
}
