package main

import (
	"github.com/archiveflax/gslauncher/internal/gui"
	"github.com/archiveflax/gslauncher/internal/settings"
)

func main() {
	settings.Load()

	go mainLoop()

	app := gui.NewApp()
	app.Run()
}
