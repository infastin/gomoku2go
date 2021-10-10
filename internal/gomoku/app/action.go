package gomoku

import (
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
)

type Action struct {
	*gio.SimpleAction
}

func NewAction(name string, parameterType *glib.VariantType, f interface{}) *Action {
	act := &Action{}
	act.SimpleAction = gio.NewSimpleAction(name, parameterType)
	act.Connect("activate", f)

	return act
}
