package main

import (
	"github.com/archiveflax/gslauncher/internal/gui"
	"github.com/archiveflax/gslauncher/internal/settings"
)

const groovestatsUrl = "http://localhost:12345" // XXX

func main() {
	settings.Load()

	go mainLoop()

	app := gui.NewApp()
	app.Run()
}
