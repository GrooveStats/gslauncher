package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

type AutoDownloadMode int

const (
	AutoDownloadOff AutoDownloadMode = iota
	AutoDownloadOnly
	AutoDownloadAndUnpack
)

func (m *AutoDownloadMode) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	switch s {
	case "off":
		*m = AutoDownloadOff
	case "download-only":
		*m = AutoDownloadOnly
	case "download-and-unpack":
		*m = AutoDownloadAndUnpack
	default:
		*m = AutoDownloadOff
	}

	return nil
}

func (m AutoDownloadMode) MarshalJSON() ([]byte, error) {
	var s string

	switch m {
	case AutoDownloadOff:
		s = "off"
	case AutoDownloadOnly:
		s = "download-only"
	case AutoDownloadAndUnpack:
		s = "download-and-unpack"
	default:
		s = "off"
	}

	return json.Marshal(s)
}

type Settings struct {
	FirstLaunch      bool `json:"-"`
	SmExePath        string
	SmDataDir        string
	AutoDownloadMode AutoDownloadMode
	UserUnlocks      bool

	// debug settings, not stored in the json
	Debug                  bool   `json:"-"`
	FakeGs                 bool   `json:"-"`
	FakeGsNetworkError     bool   `json:"-"`
	FakeGsNetworkDelay     int    `json:"-"`
	FakeGsNewSessionResult string `json:"-"`
	FakeGsSubmitResult     string `json:"-"`
	FakeGsRpg              bool   `json:"-"`
	GrooveStatsUrl         string `json:"-"`
}

var settings Settings = getDefaults()

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

	err = json.Unmarshal(data, &settings)
	if err != nil {
		return err
	}

	settings.FirstLaunch = false
	return nil
}

func Update(newSettings Settings) {
	settings = newSettings
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

	grooveStatsUrl := "https://api.groovestats.com"
	if debug {
		grooveStatsUrl = "http://localhost:9090"
	}

	return Settings{
		FirstLaunch:      true,
		SmExePath:        smExePath,
		SmDataDir:        smDataDir,
		AutoDownloadMode: AutoDownloadOff,
		UserUnlocks:      false,

		Debug:                  debug,
		FakeGs:                 debug,
		FakeGsNetworkError:     false,
		FakeGsNetworkDelay:     0,
		FakeGsNewSessionResult: "OK",
		FakeGsSubmitResult:     "score-added",
		FakeGsRpg:              true,
		GrooveStatsUrl:         grooveStatsUrl,
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
