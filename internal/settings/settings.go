package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
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
	SmSaveDir        string
	SmSongsDir       string
	SmLogsDir        string
	AutoDownloadMode AutoDownloadMode
	UserUnlocks      bool
	AutoLaunch       bool
	GrooveStatsUrl   string

	// debug settings, not stored in the json
	Debug                  bool   `json:"-"`
	FakeGs                 bool   `json:"-"`
	FakeGsNetworkError     bool   `json:"-"`
	FakeGsNetworkDelay     int    `json:"-"`
	FakeGsNewSessionResult string `json:"-"`
	FakeGsSubmitResult     string `json:"-"`
	FakeGsRpg              bool   `json:"-"`
	FakeGsItl              bool   `json:"-"`
	GrooveStatsUrl         string `json:"-"`

	// backwards compatibility fields
	SmDataDir string `json:",omitempty"`
}

var settings = Settings{
	FirstLaunch:      true,
	SmExePath:        "",
	SmSaveDir:        "",
	SmSongsDir:       "",
	SmLogsDir:        "",
	AutoDownloadMode: AutoDownloadOff,
	UserUnlocks:      false,
	AutoLaunch:       false,

	Debug:                  debug,
	FakeGs:                 false,
	FakeGsNetworkError:     false,
	FakeGsNetworkDelay:     0,
	FakeGsNewSessionResult: "OK",
	FakeGsSubmitResult:     "score-added",
	FakeGsRpg:              true,
	FakeGsItl:              true,
	GrooveStatsUrl:         "https://api.groovestats.com",
}

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

	// v1.0.0 compat
	if settings.SmDataDir != "" {
		settings.SmSaveDir = filepath.Join(settings.SmDataDir, "Save")
		settings.SmSongsDir = filepath.Join(settings.SmDataDir, "Songs")
		settings.SmLogsDir = filepath.Join(settings.SmDataDir, "Logs")
		settings.SmDataDir = ""
	}

	settings.FirstLaunch = false
	return nil
}

func Update(newSettings Settings) {
	settings = newSettings
}

func DetectSM() {
	smExePath, smSaveDir, smSongsDir, smLogsDir := detectSM()

	settings.SmExePath = smExePath
	settings.SmSaveDir = smSaveDir
	settings.SmSongsDir = smSongsDir
	settings.SmLogsDir = smLogsDir
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
