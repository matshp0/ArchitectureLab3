package painter

import (
	"github.com/matshp0/ArchitectureLab3/ui"
	"golang.org/x/exp/shiny/screen"
	"image/color"
)

// Operation змінює вхідну текстуру.
type Operation interface {
	// Do виконує зміну операції, повертаючи true, якщо текстура вважається готовою для відображення.
	Do(t screen.Texture) (ready bool)
}

// UpdateOp операція, яка не змінює текстуру, але сигналізує, що текстуру потрібно розглядати як готову.
var UpdateOp = updateOp{}

type Command struct {
	F       OperationFunc
	Options map[string]float32
}
type updateOp struct{}

// OperationFunc використовується для перетворення `функції оновлення текстури в Operation.
type OperationFunc func(t screen.Texture, options map[string]float32)

func (f updateOp) Do(t screen.Texture) bool {
	return false
}

func (f Command) Do(t screen.Texture) bool {
	f.F(t, f.Options)
	return false
}

// WhiteFill WhiteFill зафарбовує текстуру у білий колір. Може бути використана як Operation через OperationFunc(WhiteFill).
var WhiteFill OperationFunc = func(t screen.Texture, options map[string]float32) {
	t.Fill(t.Bounds(), color.White, screen.Src)
}

// GreenFill зафарбовує текстуру у зелений колір. Може бути використана як Operation через OperationFunc(GreenFill).
var GreenFill OperationFunc = func(t screen.Texture, options map[string]float32) {
	t.Fill(t.Bounds(), ui.Green, screen.Src)
}
