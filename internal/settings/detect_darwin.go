package settings

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

var locations = []string{
	"/Applications",
	"/Applications/ITGmania",
	"/Applications/StepMania-5.1.0",
	"/Applications/OutFox",
	"/Applications/StepMania",
	"/Applications/StepMania-5.0.12",
}

func detectSM() (string, string, string, string) {
	for _, installDir := range locations {
		smAppPath := filepath.Join(installDir, "StepMania.app")
		ofAppPath := filepath.Join(installDir, "OutFox.app")
		itgmAppPath := filepath.Join(installDir, "ITGmania.app")

		if _, err := os.Stat(itgmAppPath); err == nil {
			smAppPath = itgmAppPath
		} else if _, err := os.Stat(ofAppPath); err == nil {
			smAppPath = ofAppPath
		}

		_, err := os.Stat(smAppPath)
		if err != nil {
			continue
		}

		var smSaveDir string
		var smSongsDir string
		var smLogsDir string

		if installDir != "/Applications" {
			// portable installation?
			_, err = os.Stat(filepath.Join(installDir, "Portable.ini"))
			if err == nil {
				smSaveDir = filepath.Join(installDir, "Save")
				smSongsDir = filepath.Join(installDir, "Songs")
				smLogsDir = filepath.Join(installDir, "Logs")

				return smAppPath, smSaveDir, smSongsDir, smLogsDir
			}
		}

		// Query the SM version.
		smExePath := filepath.Join(smAppPath, "Contents", "MacOS", "StepMania")
		ofExePath := filepath.Join(smAppPath, "Contents", "MacOS", "OutFox")

		if _, err := os.Stat(ofExePath); err == nil {
			smExePath = ofExePath
		}

		cmd := exec.Command(smExePath, "--version")
		cmd.Dir = filepath.Dir(smExePath)

		out, err := cmd.Output()
		if err != nil {
			return smAppPath, "", "", ""
		}

		pattern := regexp.MustCompile(`(?m)^(StepMania|OutFox|ITGmania)(\d\.[\d+]+)`)
		m := pattern.FindSubmatch(out)
		if len(m) < 2 {
			return smAppPath, "", "", ""
		}
		isOutFox := string(m[1]) == "OutFox"
		isITGmania := string(m[1]) == "ITGmania"
		version := string(m[1])

		if isOutFox {
			homeDir, err := os.UserHomeDir()
			if err == nil {
				smSaveDir = filepath.Join(homeDir, "Library", "Preferences", "Project OutFox")
				smLogsDir = filepath.Join(homeDir, "Library", "Logs", "Project OutFox")
			}

			smSongsDir = filepath.Join(installDir, "Songs")
		} else if isITGmania {
			homeDir, err := os.UserHomeDir()
			if err == nil {
				smSaveDir = filepath.Join(homeDir, "Library", "Preferences", "ITGmania ")
				smSongsDir = filepath.Join(homeDir, "Library", "Application Support", "ITGmania", "Songs")
				smLogsDir = filepath.Join(homeDir, "Library", "Logs", "ITGmania")
			}
		} else {
			homeDir, err := os.UserHomeDir()
			if err == nil {
				smSaveDir = filepath.Join(homeDir, "Library", "Preferences", "StepMania "+version)
				smSongsDir = filepath.Join(homeDir, "Library", "Application Support", "StepMania "+version, "Songs")
				smLogsDir = filepath.Join(homeDir, "Library", "Logs", "StepMania "+version)
			}
		}

		return smAppPath, smSaveDir, smSongsDir, smLogsDir
	}

	return "", "", "", ""
}
