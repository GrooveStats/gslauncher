package main

import (
	"log"

	"github.com/GrooveStats/gslauncher/internal/gui"
	"github.com/GrooveStats/gslauncher/internal/settings"
	"github.com/GrooveStats/gslauncher/internal/unlocks"
)

func main() {
	settings.Load()

	unlockManager, err := unlocks.NewManager()
	if err != nil {
		log.Print("failed to initialize downloader: ", err)
		return
	}

	app := gui.NewApp(unlockManager)
	app.Run()
}
