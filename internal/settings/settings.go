package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

type Settings struct {
	SmExePath    string
	SmDataDir    string
	AutoDownload bool
	AutoUnpack   bool
	UserUnlocks  bool

	// debug settings, not stored in the json
	Debug           bool   `json:"-"`
	FakeGroovestats bool   `json:"-"`
	GrooveStatsUrl  string `json:"-"`
}

var settings Settings = getDefaults()
var updateCallback func(Settings, Settings)

func Get() Settings {
	return settings
}

func Load() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	settingsPath := filepath.Join(configDir, "groovestats-launcher", "settings.json")

	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &settings)
}

func Update(newSettings Settings) {
	oldSettings := settings
	settings = newSettings

	if updateCallback != nil {
		updateCallback(oldSettings, newSettings)
	}
}

func SetUpdateCallback(callback func(Settings, Settings)) {
	updateCallback = callback
}

func getDefaults() Settings {
	var smExePath string
	var smDataDir string

	switch runtime.GOOS {
	case "windows":
		smExePath = "C:\\Games\\StepMania 5.1\\Program\\StepMania.exe"

		configDir, err := os.UserConfigDir()
		if err == nil {
			smDataDir = filepath.Join(configDir, "StepMania 5.1")
		}
	default:
		smExePath = "/usr/local/bin/stepmania"

		homeDir, err := os.UserHomeDir()
		if err == nil {
			smDataDir = filepath.Join(homeDir, ".stepmania-5.1")
		}
	}

	grooveStatsUrl := "https://www.groovestats.com"
	if debug {
		grooveStatsUrl = "http://localhost:9090"
	}

	return Settings{
		SmExePath:    smExePath,
		SmDataDir:    smDataDir,
		AutoDownload: false,
		AutoUnpack:   false,
		UserUnlocks:  false,

		Debug:           debug,
		FakeGroovestats: false,
		GrooveStatsUrl:  grooveStatsUrl,
	}
}

func Save() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	settingsDir := filepath.Join(configDir, "groovestats-launcher")
	settingsPath := filepath.Join(settingsDir, "settings.json")

	info, err := os.Stat(settingsDir)
	if err != nil || !info.IsDir() {
		err := os.Mkdir(settingsDir, os.ModeDir|0700)
		if err != nil {
			return err
		}
	}

	data, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	return os.WriteFile(settingsPath, data, 0600)
}
