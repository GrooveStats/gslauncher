package settings

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func detectSM() (string, string, string, string) {
	cmd := exec.Command("which", "itgmania")

	out, err := cmd.Output()
	if err != nil {
		cmd = exec.Command("which", "stepmania")

		out, err = cmd.Output()
		if err != nil {
			return "", "", "", ""
		}
	}
	smExePath := strings.TrimSpace(string(out))

	// follow if it's a symlink
	target, err := os.Readlink(smExePath)
	if err == nil {
		if !strings.HasPrefix(target, "/") {
			target = filepath.Join(filepath.Dir(smExePath), target)
		}
		smExePath = target
	}

	// portable installation?
	installDir := filepath.Dir(smExePath)
	_, err = os.Stat(filepath.Join(installDir, "portable.ini"))
	if err == nil {
		smSaveDir := filepath.Join(installDir, "Save")
		smSongsDir := filepath.Join(installDir, "Songs")
		smLogsDir := filepath.Join(installDir, "Logs")
		return smExePath, smSaveDir, smSongsDir, smLogsDir
	}

	// Query the SM version. We also have to set the working directory,
	// because SM 5.3 (outfox) for Linux searches for bundled shared
	// libraries in the current working directory.
	cmd = exec.Command(smExePath, "--version")
	cmd.Dir = filepath.Dir(smExePath)

	out, err = cmd.Output()
	if err != nil {
		return smExePath, "", "", ""
	}

	pattern := regexp.MustCompile(`(?m)^(StepMania|OutFox|ITGmania)(\d\.[\d+]+)`)
	m := pattern.FindSubmatch(out)
	if len(m) < 2 {
		return smExePath, "", "", ""
	}
	isOutFox := string(m[1]) == "OutFox"
	isITGmania := string(m[1]) == "ITGmania"
	version := string(m[2])

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return smExePath, "", "", ""
	}

	smDataDir := filepath.Join(homeDir, ".stepmania-"+version)
	if isOutFox {
		smDataDir = filepath.Join(homeDir, ".project-outfox")
	} else if isITGmania {
		smDataDir = filepath.Join(homeDir, ".itgmania")
	}

	smSaveDir := filepath.Join(smDataDir, "Save")
	smSongsDir := filepath.Join(smDataDir, "Songs")
	smLogsDir := filepath.Join(smDataDir, "Logs")

	return smExePath, smSaveDir, smSongsDir, smLogsDir
}
