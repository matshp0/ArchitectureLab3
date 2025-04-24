package painter

import (
	"github.com/matshp0/ArchitectureLab3/ui"
	"golang.org/x/exp/shiny/screen"
	"image"
	"image/color"
)

type Operation interface {
	Do(t screen.Texture) bool
}

type OperationFunc func(t screen.Texture, options map[string]float32)

type Command struct {
	F       OperationFunc
	Options map[string]float32
}

type updateOp struct{}

func (f updateOp) Do(t screen.Texture) bool {
	return false
}

func (f Command) Do(t screen.Texture) bool {
	f.F(t, f.Options)
	return false
}

var UpdateOp = updateOp{}

type Figure struct {
	RelX, RelY float32
}

type State struct {
	Figures         []Figure
	BGRect          *image.Rectangle
	BackgroundColor color.Color
}

var GlobalState = &State{
	BackgroundColor: color.Black,
}

func relToAbs(v float32, max int) int {
	return int(v * float32(max))
}

func drawFigures(t screen.Texture) {
	b := t.Bounds()
	for _, f := range GlobalState.Figures {
		x := relToAbs(f.RelX, b.Dx()) + b.Min.X
		y := relToAbs(f.RelY, b.Dy()) + b.Min.Y

		t.Fill(image.Rect(x-100, y-50, x, y+50), ui.Blue, screen.Src)
		t.Fill(image.Rect(x, y-150, x+100, y+150), ui.Blue, screen.Src)
	}
}

func drawBGRect(t screen.Texture) {
	if GlobalState.BGRect != nil {
		t.Fill(*GlobalState.BGRect, color.Black, screen.Src)
	}
}

var WhiteFill OperationFunc = func(t screen.Texture, options map[string]float32) {
	t.Fill(t.Bounds(), color.White, screen.Src)
	GlobalState.BackgroundColor = color.White
	drawBGRect(t)
	drawFigures(t)
}

var GreenFill OperationFunc = func(t screen.Texture, options map[string]float32) {
	t.Fill(t.Bounds(), ui.Green, screen.Src)
	GlobalState.BackgroundColor = ui.Green
	drawBGRect(t)
	drawFigures(t)
}

var BGRect OperationFunc = func(t screen.Texture, options map[string]float32) {
	b := t.Bounds()
	x1 := relToAbs(options["x1"], b.Dx()) + b.Min.X
	y1 := relToAbs(options["y1"], b.Dy()) + b.Min.Y
	x2 := relToAbs(options["x2"], b.Dx()) + b.Min.X
	y2 := relToAbs(options["y2"], b.Dy()) + b.Min.Y

	r := image.Rect(x1, y1, x2, y2).Canon()
	GlobalState.BGRect = &r
	t.Fill(r, color.Black, screen.Src)
	drawFigures(t)
}

var Figure1 OperationFunc = func(t screen.Texture, options map[string]float32) {
	GlobalState.Figures = append(GlobalState.Figures, Figure{
		RelX: options["x"],
		RelY: options["y"],
	})

	drawFigures(t)
}

var Move OperationFunc = func(t screen.Texture, options map[string]float32) {
	dx := options["x"]
	dy := options["y"]

	for i := range GlobalState.Figures {
		GlobalState.Figures[i].RelX = dx
		GlobalState.Figures[i].RelY = dy
	}

	t.Fill(t.Bounds(), GlobalState.BackgroundColor, screen.Src)
	drawBGRect(t)
	drawFigures(t)
}

var Reset OperationFunc = func(t screen.Texture, options map[string]float32) {
	GlobalState.Figures = nil
	GlobalState.BGRect = nil
	t.Fill(t.Bounds(), GlobalState.BackgroundColor, screen.Src)
}
