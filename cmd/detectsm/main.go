package main

import (
	"fmt"
	"github.com/GrooveStats/gslauncher/internal/settings"
)

func main() {
	settings.DetectSM()

	data := settings.Get()
	fmt.Printf("exe: %s\n", data.SmExePath)
	fmt.Printf("Save/: %s\n", data.SmSaveDir)
	fmt.Printf("Songs/: %s\n", data.SmSongsDir)
	fmt.Printf("Logs/: %s\n", data.SmLogsDir)
}
