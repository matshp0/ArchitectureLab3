package main

import (
	"net/http"

	"github.com/matshp0/ArchitectureLab3/painter"
	"github.com/matshp0/ArchitectureLab3/painter/lang"
	"github.com/matshp0/ArchitectureLab3/ui"
)

func main() {
	var (
		pv ui.Visualizer // Візуалізатор створює вікно та малює у ньому.

		// Потрібні для частини 2.
		opLoop = *painter.NewLoop()
		parser lang.Parser // Парсер команд.
	)

	pv.Debug = false
	pv.Title = "Simple painter"

	pv.OnScreenReady = opLoop.Start
	opLoop.Receiver = &pv

	go func() {
		http.Handle("/", lang.HttpHandler(&opLoop, &parser))
		_ = http.ListenAndServe("localhost:17000", nil)
	}()

	pv.Main()
	opLoop.StopAndWait()
}
