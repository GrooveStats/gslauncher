package main

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/GrooveStats/gslauncher/internal/gui"
	"github.com/GrooveStats/gslauncher/internal/settings"
	"github.com/GrooveStats/gslauncher/internal/unlocks"
)

func redirectLog() {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Print("failed to get cache directory: ", err)
		return
	}

	filename := filepath.Join(cacheDir, "groovestats-launcher", "log.txt")

	logFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Print("failed to open log file: ", err)
		return
	}

	if settings.Get().Debug {
		log.SetOutput(io.MultiWriter(os.Stderr, logFile))
	} else {
		log.SetOutput(logFile)
	}
}

func main() {
	redirectLog()

	settings.Load()

	unlockManager, err := unlocks.NewManager()
	if err != nil {
		log.Print("failed to initialize downloader: ", err)
		return
	}

	app := gui.NewApp(unlockManager)
	app.Run()
}
