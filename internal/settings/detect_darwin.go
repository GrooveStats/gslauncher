package settings

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

var locations = []string{
	"/Applications",
	"/Applications/StepMania-5.1.0",
	"/Applications/StepMania",
	"/Applications/StepMania-5.0.12",
}

func detectSM() (string, string, string, string) {
	for _, installDir := range locations {
		smAppPath := filepath.Join(installDir, "StepMania.app")

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
		cmd := exec.Command(smExePath, "--version")
		cmd.Dir = filepath.Dir(smExePath)

		out, err := cmd.Output()
		if err != nil {
			return smAppPath, "", "", ""
		}

		pattern := regexp.MustCompile(`^StepMania(5\.[\d+]+)`)
		m := pattern.FindSubmatch(out)
		if len(m) < 2 {
			return smAppPath, "", "", ""
		}
		version := string(m[1])

		switch version {
		case "5.0", "5.1":
			homeDir, err := os.UserHomeDir()
			if err == nil {
				smSaveDir = filepath.Join(homeDir, "Library", "Preferences", "StepMania "+version)
				smSongsDir = filepath.Join(homeDir, "Library", "Application Support", "StepMania "+version, "Songs")
				smLogsDir = filepath.Join(homeDir, "Library", "Logs", "StepMania "+version)
			}
		case "5.3":
			homeDir, err := os.UserHomeDir()
			if err == nil {
				smSaveDir = filepath.Join(homeDir, "Library", "Preferences", "StepMania 5.3")
				smLogsDir = filepath.Join(homeDir, "Library", "Logs", "StepMania 5.3")
			}

			smSongsDir = filepath.Join(installDir, "Songs")
		}

		return smAppPath, smSaveDir, smSongsDir, smLogsDir
	}

	return "", "", "", ""
}
