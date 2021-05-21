package main

import (
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/GrooveStats/gslauncher/internal/gui"
	"github.com/GrooveStats/gslauncher/internal/settings"
	"github.com/GrooveStats/gslauncher/internal/unlocks"
	"github.com/GrooveStats/gslauncher/internal/version"
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
	autolaunch := flag.Bool("autolaunch", false, "automatically launch StepMania")
	flag.Parse()

	redirectLog()
	log.Printf("GrooveStats Launcher %s (%s %s)", version.Formatted(), runtime.GOOS, runtime.GOARCH)

	settings.Load()
	if settings.Get().FirstLaunch {
		settings.DetectSM()
	}

	unlockManager, err := unlocks.NewManager()
	if err != nil {
		log.Print("failed to initialize downloader: ", err)
		return
	}

	app := gui.NewApp(unlockManager, *autolaunch)
	app.Run()
}
