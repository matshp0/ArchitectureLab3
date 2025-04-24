package ui

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/imageutil"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

var (
	Green      = color.RGBA{R: 140, G: 198, B: 158, A: 255}
	Blue       = color.RGBA{R: 11, G: 132, B: 253, A: 255}
	WindowSize = 800
)

type Visualizer struct {
	Title         string
	Debug         bool
	OnScreenReady func(s screen.Screen)

	w    screen.Window
	tx   chan screen.Texture
	done chan struct{}

	sz  size.Event
	pos image.Point
}

func (pw *Visualizer) Main() {
	pw.tx = make(chan screen.Texture)
	pw.done = make(chan struct{})
	driver.Main(pw.run)
}

func (pw *Visualizer) Update(t screen.Texture) {
	pw.tx <- t
}

func (pw *Visualizer) run(s screen.Screen) {
	w, err := s.NewWindow(&screen.NewWindowOptions{
		Width:  WindowSize,
		Height: WindowSize,
		Title:  pw.Title,
	})
	pw.pos = image.Point{
		X: WindowSize / 2,
		Y: WindowSize / 2,
	}
	if err != nil {
		log.Fatal("Failed to initialize the app window:", err)
	}
	defer func() {
		w.Release()
		close(pw.done)
	}()

	if pw.OnScreenReady != nil {
		pw.OnScreenReady(s)
	}

	pw.w = w

	events := make(chan any)
	go func() {
		for {
			e := w.NextEvent()
			if pw.Debug {
				log.Printf("new event: %v", e)
			}
			if detectTerminate(e) {
				close(events)
				break
			}
			events <- e
		}
	}()

	var t screen.Texture

	for {
		select {
		case e, ok := <-events:
			if !ok {
				return
			}
			pw.handleEvent(e, t)

		case t = <-pw.tx:
			w.Send(paint.Event{})
		}
	}
}

func detectTerminate(e any) bool {
	switch e := e.(type) {
	case lifecycle.Event:
		if e.To == lifecycle.StageDead {
			return true // Window destroy initiated.
		}
	case key.Event:
		if e.Code == key.CodeEscape {
			return true // Esc pressed.
		}
	}
	return false
}

func (pw *Visualizer) handleEvent(e any, t screen.Texture) {
	switch e := e.(type) {

	case size.Event: // Оновлення даних про розмір вікна.
		pw.sz = e
		fmt.Println(e)

	case error:
		log.Printf("ERROR: %s", e)

	case mouse.Event:
		if t == nil {
			if e.Button == mouse.ButtonLeft && e.Direction == mouse.DirPress {
				x := int(e.X)
				y := int(e.Y)
				pw.pos = image.Point{X: x, Y: y}
				pw.w.Send(paint.Event{})

			}
		}

	case paint.Event:
		// Малювання контенту вікна.
		if t == nil {
			pw.drawDefaultUI()
		} else {
			pw.w.Scale(pw.sz.Bounds(), t, t.Bounds(), draw.Src, nil)
		}
		pw.w.Publish()
	}
}

func (pw *Visualizer) drawDefaultUI() {
	bounds := pw.sz.Bounds()
	pw.w.Fill(bounds, Green, draw.Src) // Фон.

	pw.drawTShape()
	// Малювання білої рамки.
	for _, br := range imageutil.Border(bounds, 10) {
		pw.w.Fill(br, color.White, draw.Src)
	}
}

func (pw *Visualizer) drawTShape() {
	pw.w.Fill(image.Rect(pw.pos.X-100, pw.pos.Y-50, pw.pos.X, pw.pos.Y+50), Blue, 0)
	pw.w.Fill(image.Rect(pw.pos.X, pw.pos.Y-150, pw.pos.X+100, pw.pos.Y+150), Blue, 0)
}
