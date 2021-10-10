package view

import (
	_ "embed"
	"fmt"
	"math"

	"github.com/diamondburned/gotk4/pkg/cairo"
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type Shape uint

const (
	Circle Shape = iota
	Cross
)

const (
	strokeCoef = 32
)

type Field struct {
	X, Y uint
	Sh   Shape
}

type BoardArea struct {
	*gtk.DrawingArea

	clickHandler  func(x, y uint)
	redrawHandler func()

	sneaky  *gtk.Button
	surface *cairo.Surface
	press   *gtk.GestureClick
	cells   uint
}

func newBoardArea(builder *gtk.Builder) *BoardArea {
	board := &BoardArea{}

	board.sneaky = builder.GetObject("sneaky").Cast().(*gtk.Button)
	board.sneaky.GrabFocus()

	board.DrawingArea = builder.GetObject("board").Cast().(*gtk.DrawingArea)
	board.SetDrawFunc(board.drawFunc)
	board.ConnectAfter("resize", board.onResize)

	return board
}

func (board *BoardArea) Init(cells uint) {
	board.cells = cells

	board.press = gtk.NewGestureClick()
	board.press.SetButton(gdk.BUTTON_PRIMARY)
	board.press.ConnectPressed(board.onPress)
	board.AddController(board.press)
}

func (board *BoardArea) onResize(_ *gtk.DrawingArea, width, height int) {
	if surf := board.GetNative().Surface(); surf != nil {
		board.surface = surf.CreateSimilarSurface(cairo.CONTENT_COLOR, width, height)

		if board.cells != 0 {
			board.Draw()

			if board.redrawHandler != nil {
				board.redrawHandler()
			}
		}
	}
}

func (board *BoardArea) onPress(nPress int, x, y float64) {
	board.sneaky.GrabFocus()

	if board.clickHandler == nil {
		return
	}

	width := board.Width()
	height := board.Height()

	min := math.Min(float64(width), float64(height))
	size := (float64(min) * 2) / 3

	fcells := float64(board.cells)
	linew := float64(size / (strokeCoef * fcells))
	csize := size / fcells

	x0 := float64(width)/2 - (size / 2) + linew
	y0 := float64(height)/2 - (size / 2) + linew

	x1 := float64(width)/2 + (size / 2)
	y1 := float64(height)/2 + (size / 2)

	if (x < x0 || x > x1) || (y < y0 || y > y1) {
		return
	}

	tx := x - x0
	ty := y - y0

	cellx := math.Floor(tx / csize)
	celly := math.Floor(ty / csize)

	if tx-(csize*cellx) > csize-linew || ty-(csize*celly) > csize-linew {
		return
	}

	board.clickHandler(uint(cellx), uint(celly))
}

func (board *BoardArea) drawFunc(_ *gtk.DrawingArea, cr *cairo.Context, width, height int) {
	cr.SetSourceSurface(board.surface, 0, 0)
	cr.Paint()
}

func (board *BoardArea) Draw() {
	width := board.Width()
	height := board.Height()
	min := math.Min(float64(width), float64(height))
	size := (float64(min) * 2) / 3

	fcells := float64(board.cells)
	linew := float64(size / (strokeCoef * fcells))

	sctx := board.StyleContext()
	fg, _ := sctx.LookupColor("theme_fg_color")
	bg, _ := sctx.LookupColor("theme_bg_color")

	cr := cairo.Create(board.surface)
	cr.SetSourceRGBA(float64(bg.Red()),
		float64(bg.Green()),
		float64(bg.Blue()),
		float64(bg.Alpha()))
	cr.Paint()

	cr.SetSourceRGBA(float64(fg.Red()),
		float64(fg.Green()),
		float64(fg.Blue()),
		float64(fg.Alpha()))

	cr.Rectangle(float64(width)/2-(size/2), float64(height)/2-(size/2), size, size)
	cr.SetLineWidth(linew)
	cr.SetLineJoin(cairo.LINE_JOIN_MITER)

	cr.Translate(float64(width)/2-(size/2), float64(height)/2-(size/2))

	for i := uint(1); i < board.cells; i++ {
		fi := float64(i)

		cr.MoveTo((size/fcells)*fi, 0)
		cr.LineTo((size/fcells)*fi, size)

		cr.MoveTo(0, (size/fcells)*fi)
		cr.LineTo(size, (size/fcells)*fi)
	}

	cr.Stroke()

	board.QueueDraw()
}

func (board *BoardArea) drawCircle(x, y uint, queue bool) error {
	if board.cells == 0 {
		return fmt.Errorf("boardArea hasn't been initialized")
	}

	width := board.Width()
	height := board.Height()
	min := math.Min(float64(width), float64(height))
	size := (float64(min) * 2) / 3

	fcells := float64(board.cells)
	csize := size / fcells
	radius := csize / 3
	linew := float64(size / (strokeCoef * fcells))

	cx := csize*float64(x) + csize/2
	cy := csize*float64(y) + csize/2

	sctx := board.StyleContext()
	fg, _ := sctx.LookupColor("theme_fg_color")

	cr := cairo.Create(board.surface)
	cr.Translate(float64(width)/2-(size/2), float64(height)/2-(size/2))

	cr.SetSourceRGBA(float64(fg.Red()),
		float64(fg.Green()),
		float64(fg.Blue()),
		float64(fg.Alpha()))

	cr.Arc(cx, cy, radius, 0, 2*math.Pi)
	cr.SetLineWidth(linew)
	cr.Stroke()

	if queue {
		board.QueueDraw()
	}

	return nil
}

func (board *BoardArea) DrawCircle(x, y uint) error {
	return board.drawCircle(x, y, true)
}

func (board *BoardArea) drawCross(x, y uint, queue bool) error {
	if board.cells == 0 {
		return fmt.Errorf("boardArea hasn't been initialized")
	}

	width := board.Width()
	height := board.Height()
	min := math.Min(float64(width), float64(height))
	size := (float64(min) * 2) / 3

	fcells := float64(board.cells)
	csize := size / fcells
	linew := float64(size / (strokeCoef * fcells))

	x0 := csize*float64(x) + csize/6
	y0 := csize*float64(y) + csize/6
	x1 := csize*float64(x) + (csize*5)/6
	y1 := csize*float64(y) + (csize*5)/6

	sctx := board.StyleContext()
	fg, _ := sctx.LookupColor("theme_fg_color")

	cr := cairo.Create(board.surface)
	cr.Translate(float64(width)/2-(size/2), float64(height)/2-(size/2))

	cr.SetSourceRGBA(float64(fg.Red()),
		float64(fg.Green()),
		float64(fg.Blue()),
		float64(fg.Alpha()))

	cr.MoveTo(x0, y0)
	cr.LineTo(x1, y1)

	cr.MoveTo(x1, y0)
	cr.LineTo(x0, y1)

	cr.SetLineWidth(linew)
	cr.Stroke()

	if queue {
		board.QueueDraw()
	}

	return nil
}

func (board *BoardArea) DrawCross(x, y uint) error {
	return board.drawCross(x, y, true)
}

func (board *BoardArea) DrawShapes(fields []Field) {
	for _, f := range fields {
		switch f.Sh {
		case Circle:
			board.drawCircle(f.X, f.Y, false)
		case Cross:
			board.drawCross(f.X, f.Y, false)
		}
	}

	board.QueueDraw()
}

func (board *BoardArea) DrawStrike(x0, y0, x1, y1 uint) error {
	if board.cells == 0 {
		return fmt.Errorf("boardArea hasn't been initialized")
	}

	width := board.Width()
	height := board.Height()

	min := math.Min(float64(width), float64(height))
	size := (float64(min) * 2) / 3

	fcells := float64(board.cells)
	linew := float64((size * 2) / (strokeCoef * fcells))
	csize := size / fcells

	var bx0 float64
	var by0 float64
	var bx1 float64
	var by1 float64

	switch {
	case x0 == x1:
		bx0 = csize*float64(x0) + csize/2
		by0 = csize*float64(y0) + csize/6
		bx1 = csize*float64(x1) + csize/2
		by1 = csize*float64(y1) + (5*csize)/6
	case y0 == y1:
		bx0 = csize*float64(x0) + csize/6
		by0 = csize*float64(y0) + csize/2
		bx1 = csize*float64(x1) + (5*csize)/6
		by1 = csize*float64(y1) + csize/2
	case y0 < y1:
		bx0 = csize*float64(x0) + csize/6
		by0 = csize*float64(y0) + csize/6
		bx1 = csize*float64(x1) + (5*csize)/6
		by1 = csize*float64(y1) + (5*csize)/6
	default:
		bx0 = csize*float64(x0) + csize/6
		by0 = csize*float64(y0) + (5*csize)/6
		bx1 = csize*float64(x1) + (5*csize)/6
		by1 = csize*float64(y1) + csize/6
	}

	sctx := board.StyleContext()
	sbg, _ := sctx.LookupColor("theme_selected_bg_color")

	cr := cairo.Create(board.surface)
	cr.Translate(float64(width)/2-(size/2), float64(height)/2-(size/2))

	cr.SetSourceRGBA(float64(sbg.Red()),
		float64(sbg.Green()),
		float64(sbg.Blue()),
		float64(sbg.Alpha()))

	cr.MoveTo(bx0, by0)
	cr.LineTo(bx1, by1)

	cr.SetLineWidth(linew)
	cr.Stroke()

	board.QueueDraw()

	return nil
}

func (board *BoardArea) Clear() {
	if board.clickHandler != nil {
		board.clickHandler = nil
	}

	if board.redrawHandler != nil {
		board.redrawHandler = nil
	}

	board.RemoveController(board.press)
	board.cells = 0

	sctx := board.StyleContext()
	bg, _ := sctx.LookupColor("theme_bg_color")

	cr := cairo.Create(board.surface)
	cr.SetSourceRGBA(float64(bg.Red()),
		float64(bg.Green()),
		float64(bg.Blue()),
		float64(bg.Alpha()))
	cr.Paint()

	board.QueueDraw()
}

func (board *BoardArea) ConnectClick(handler func(x, y uint)) error {
	if board.cells == 0 {
		return fmt.Errorf("boardArea hasn't been initialized")
	}

	if board.clickHandler == nil {
		board.clickHandler = handler
	}

	return nil
}

func (board *BoardArea) ConnectRedraw(handler func()) error {
	if board.cells == 0 {
		return fmt.Errorf("boardArea hasn't been initialized")
	}

	if board.redrawHandler == nil {
		board.redrawHandler = handler
	}

	return nil
}
