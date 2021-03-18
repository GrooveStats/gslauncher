package main

import (
	"log"

	"github.com/archiveflax/gslauncher/internal/gui"
	"github.com/archiveflax/gslauncher/internal/settings"
	"github.com/archiveflax/gslauncher/internal/unlocks"
)

func main() {
	settings.Load()

	unlockManager, err := unlocks.NewManager()
	if err != nil {
		log.Print("failed to initialize downloader: ", err)
		return
	}

	go mainLoop(unlockManager)

	app := gui.NewApp(unlockManager)
	app.Run()
}
