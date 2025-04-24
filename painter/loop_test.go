package painter

import (
	"github.com/matshp0/ArchitectureLab3/ui"
	"image"
	"image/color"
	"image/draw"
	"reflect"
	"testing"
	"time"

	"golang.org/x/exp/shiny/screen"
)

func newCommand(f OperationFunc) Command {
	return Command{
		F:       f,
		Options: nil,
	}
}

func TestLoop_ConcurrentOperations(t *testing.T) {
	l := NewLoop()
	tr := testReceiver{}
	l.Receiver = &tr

	l.Start(mockScreen{})

	var opsExecuted int
	for i := 0; i < 10; i++ {
		go func() {
			l.Post(newCommand(func(t screen.Texture, o map[string]float32) {
				opsExecuted++
			}))
		}()
	}

	time.Sleep(100 * time.Millisecond)
	l.StopAndWait()

	if opsExecuted != 10 {
		t.Errorf("Expected 10 operations executed, got %d", opsExecuted)
	}
}

func TestLoop_UpdateBeforeOperations(t *testing.T) {
	l := NewLoop()
	tr := testReceiver{}
	l.Receiver = &tr

	l.Start(mockScreen{})
	l.Post(UpdateOp) // Update before any operations
	l.Post(logOp(t, "do white fill", WhiteFill))
	l.StopAndWait()

	if tr.lastTexture == nil {
		t.Fatal("Texture should exist even after empty update")
	}
}

func TestLoop_MultipleUpdates(t *testing.T) {
	l := NewLoop()
	tr := testReceiver{}
	l.Receiver = &tr

	l.Start(mockScreen{})
	l.Post(logOp(t, "do white fill", WhiteFill))
	l.Post(UpdateOp)
	l.Post(logOp(t, "do green fill", GreenFill))
	l.Post(UpdateOp)
	l.StopAndWait()

	mt, ok := tr.lastTexture.(*mockTexture)
	if !ok {
		t.Fatal("Unexpected texture type")
	}

	if len(mt.Colors) != 2 {
		t.Errorf("Expected 2 colors, got %d", len(mt.Colors))
	}
	if mt.Colors[0] != color.White {
		t.Error("First color should be white")
	}
	if mt.Colors[1] != ui.Green {
		t.Error("Second color should be green")
	}
}

func TestLoop_EmptyQueue(t *testing.T) {
	l := NewLoop()
	tr := testReceiver{}
	l.Receiver = &tr

	l.Start(mockScreen{})
	// No operations posted
	l.StopAndWait()

	if tr.lastTexture != nil {
		t.Error("Texture should not be updated with empty queue")
	}
}

func TestLoop_OperationOrder(t *testing.T) {
	l := NewLoop()
	tr := testReceiver{}
	l.Receiver = &tr

	var executionOrder []int

	l.Start(mockScreen{})
	l.Post(newCommand(func(t screen.Texture, o map[string]float32) {
		executionOrder = append(executionOrder, 1)
	}))
	l.Post(newCommand(func(t screen.Texture, o map[string]float32) {
		executionOrder = append(executionOrder, 2)
	}))
	l.Post(UpdateOp)
	l.StopAndWait()

	expected := []int{1, 2}
	if !reflect.DeepEqual(executionOrder, expected) {
		t.Errorf("Operations executed out of order, expected %v got %v", expected, executionOrder)
	}
}

func TestLoop_Post(t *testing.T) {
	var (
		l  = NewLoop()
		tr testReceiver
	)
	l.Receiver = &tr

	var testOps []string

	l.Start(mockScreen{})
	l.Post(logOp(t, "do white fill", WhiteFill))
	l.Post(logOp(t, "do green fill", GreenFill))
	l.Post(UpdateOp)

	for i := 0; i < 3; i++ {
		go l.Post(logOp(t, "do green fill", GreenFill))
	}

	l.Post(newCommand(func(screen.Texture, map[string]float32) {
		testOps = append(testOps, "op 1")
	}))
	l.Post(newCommand(func(screen.Texture, map[string]float32) {
		testOps = append(testOps, "op 3")
	}))

	l.StopAndWait()

	if tr.lastTexture == nil {
		t.Fatal("Texture was not updated")
	}
	mt, ok := tr.lastTexture.(*mockTexture)
	if !ok {
		t.Fatal("Unexpected texture", tr.lastTexture)
	}
	if mt.Colors[0] != color.White {
		t.Error("First color is not white:", mt.Colors)
	}
	if len(mt.Colors) != 2 {
		t.Error("Unexpected size of colors:", mt.Colors)
	}

	if !reflect.DeepEqual(testOps, []string{"op 1", "op 3"}) {
		t.Error("Bad order:", testOps)
	}
}

func logOp(t *testing.T, msg string, op OperationFunc) Command {
	return newCommand(func(tx screen.Texture, options map[string]float32) {
		t.Log(msg)
		op(tx, options)
	})
}

type testReceiver struct {
	lastTexture screen.Texture
}

func (tr *testReceiver) Update(t screen.Texture) {
	tr.lastTexture = t
}

type mockScreen struct{}

func (m mockScreen) NewBuffer(size image.Point) (screen.Buffer, error) {
	panic("implement me")
}

func (m mockScreen) NewTexture(size image.Point) (screen.Texture, error) {
	return new(mockTexture), nil
}

func (m mockScreen) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	panic("implement me")
}

type mockTexture struct {
	Colors []color.Color
}

func (m *mockTexture) Release() {}

func (m *mockTexture) Size() image.Point { return size }

func (m *mockTexture) Bounds() image.Rectangle {
	return image.Rectangle{Max: m.Size()}
}

func (m *mockTexture) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {}
func (m *mockTexture) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	m.Colors = append(m.Colors, src)
}
