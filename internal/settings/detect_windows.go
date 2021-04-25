package settings

import (
	"os"
	"path/filepath"
)

const (
	sm50ExePath = "C:\\Games\\StepMania 5\\Program\\StepMania.exe"
	sm51ExePath = "C:\\Games\\StepMania 5.1\\Program\\StepMania.exe"
	sm53ExePath = "C:\\Games\\StepMania 5.3 Outfox\\Program\\StepMania.exe"
)

func detectSM() (string, string) {
	var sm50, sm51, sm53 bool

	_, err := os.Stat(sm50ExePath)
	if err == nil {
		sm50 = true
	}

	_, err = os.Stat(sm51ExePath)
	if err == nil {
		sm51 = true
	}

	_, err = os.Stat(sm53ExePath)
	if err == nil {
		sm53 = true
	}

	var smExePath, smDataDir string

	switch {
	case sm51:
		smExePath = sm51ExePath

		configDir, err := os.UserConfigDir()
		if err == nil {
			smDataDir = filepath.Join(configDir, "StepMania 5.1")
		}
	case sm53:
		// 5.3 is a portable installation by default
		smExePath = sm53ExePath
		smDataDir = "C:\\Games\\StepMania 5.3 Outfox"
	case sm50:
		smExePath = sm50ExePath

		configDir, err := os.UserConfigDir()
		if err == nil {
			smDataDir = filepath.Join(configDir, "StepMania 5")
		}
	}

	return smExePath, smDataDir
}
