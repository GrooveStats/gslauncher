package settings

import (
	"os"
	"path/filepath"
)

var dirnames = []string{
	"ITGmania",
	"StepMania 5.1",
	"Project OutFox",
	"StepMania 5.3 Outfox",
	"StepMania 5",
}

func detectSM() (string, string, string, string) {
	for _, dirname := range dirnames {
		installDir := filepath.Join("C:\\Games", dirname)
		smExePath := filepath.Join(installDir, "Program\\StepMania.exe")
		ofExePath := filepath.Join(installDir, "Program\\OutFox.exe")
		itgmExePath := filepath.Join(installDir, "Program\\ITGmania.exe")

		if _, err := os.Stat(itgmExePath); err == nil {
			smExePath = itgmExePath
		} else if _, err := os.Stat(ofExePath); err == nil {
			smExePath = ofExePath
		}

		_, err := os.Stat(smExePath)
		if err != nil {
			continue
		}

		var smDataDir string
		var smSaveDir string
		var smSongsDir string
		var smLogsDir string

		// portable installation?
		_, err = os.Stat(filepath.Join(installDir, "portable.ini"))
		if err == nil {
			smDataDir = installDir
		} else {
			configDir, err := os.UserConfigDir()
			if err == nil {
				smDataDir = filepath.Join(configDir, dirname)
			}
		}

		if smDataDir != "" {
			smSaveDir = filepath.Join(smDataDir, "Save")
			smSongsDir = filepath.Join(smDataDir, "Songs")
			smLogsDir = filepath.Join(smDataDir, "Logs")
		}

		return smExePath, smSaveDir, smSongsDir, smLogsDir
	}

	return "", "", "", ""
}
