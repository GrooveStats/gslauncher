package main

import (
	"bytes"
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

const MAX_LOG_SIZE = 1024 * 1024 // 1 MiB

func redirectLog(cacheDir string) {
	filename := filepath.Join(cacheDir, "groovestats-launcher", "log.txt")

	old, err := os.ReadFile(filename)
	if err == nil {
		if len(old) > MAX_LOG_SIZE {
			old = old[len(old)-MAX_LOG_SIZE:]
			idx := bytes.IndexByte(old, byte('\n'))
			old = old[idx+1:]
		}
	}

	logFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Print("failed to open log file: ", err)
		return
	}

	if old != nil {
		logFile.Write(old)
		logFile.WriteString("-----\n")
	}

	if settings.Get().Debug {
		log.SetOutput(io.MultiWriter(os.Stderr, logFile))
	} else {
		log.SetOutput(logFile)
	}
}

func main() {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Print("failed to get cache directory: ", err)
		return
	}

	autolaunch := flag.Bool("autolaunch", false, "automatically launch StepMania")
	cacheDir := flag.String("cachedir", userCacheDir, "set the cache location")
	flag.Parse()

	redirectLog(*cacheDir)
	log.Printf("GrooveStats Launcher %s (%s %s)", version.Formatted(), runtime.GOOS, runtime.GOARCH)

	settings.Load()
	if settings.Get().FirstLaunch {
		settings.DetectSM()
	}

	unlockManager, err := unlocks.NewManager(*cacheDir)
	if err != nil {
		log.Print("failed to initialize downloader: ", err)
		return
	}

	app := gui.NewApp(unlockManager, *autolaunch, *cacheDir)
	app.Run()
}
