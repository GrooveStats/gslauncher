package unlocks

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/GrooveStats/gslauncher/internal/settings"
)

type DownloadStatus int

const (
	NotDownloaded DownloadStatus = iota
	Downloading
	Downloaded
)

type UnpackStatus int

const (
	NotUnpacked UnpackStatus = iota
	Unpacking
	Unpacked
)

type UserData struct {
	ProfileName  string
	UnpackStatus UnpackStatus
	UnpackError  error
}

type Unlock struct {
	DownloadUrl      string
	RpgName          string
	QuestTitle       string
	SongDescriptions []string
	DownloadStatus   DownloadStatus
	DownloadError    error
	DownloadSize     int64
	DownloadProgress int64
	Users            []*UserData

	queue chan interface{}
}

type Manager struct {
	DownloadDir string
	Unlocks     []*Unlock

	updateCallback func(*Unlock)
}

func NewManager() (*Manager, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}

	downloadDir := filepath.Join(cacheDir, "groovestats-launcher", "unlocks")

	err = os.MkdirAll(downloadDir, os.ModeDir|0700)
	if err != nil {
		return nil, err
	}

	manager := Manager{
		DownloadDir: downloadDir,
		Unlocks:     make([]*Unlock, 0),
	}

	return &manager, nil
}

func (manager *Manager) AddUnlock(questTitle, url, rpgName, profileName string, songDescriptions []string) {
	for _, unlock := range manager.Unlocks {
		if unlock.DownloadUrl == url {
			user := &UserData{
				ProfileName: profileName,
			}
			unlock.Users = append(unlock.Users, user)
			manager.detectUnpackStatus(unlock, user)

			manager.updateCallback(unlock)

			mode := settings.Get().AutoDownloadMode
			if mode == settings.AutoDownloadAndUnpack {
				unlock.QueueUnpack(user)
			}
			return
		}
	}

	user := &UserData{
		ProfileName: profileName,
	}

	unlock := &Unlock{
		DownloadUrl:      url,
		RpgName:          rpgName,
		QuestTitle:       questTitle,
		SongDescriptions: songDescriptions,
		Users: []*UserData{
			user,
		},

		queue: make(chan interface{}, 10),
	}
	manager.detectDownloadStatus(unlock)
	manager.detectUnpackStatus(unlock, unlock.Users[0])

	manager.Unlocks = append(manager.Unlocks, unlock)

	manager.updateCallback(unlock)

	go manager.processQueue(unlock)

	mode := settings.Get().AutoDownloadMode
	if mode == settings.AutoDownloadOnly || mode == settings.AutoDownloadAndUnpack {
		unlock.QueueDownload()
		if mode == settings.AutoDownloadAndUnpack {
			unlock.QueueUnpack(user)
		}
	}
}

func (manager *Manager) SetUpdateCallback(callback func(*Unlock)) {
	manager.updateCallback = callback
}

func (manager *Manager) detectDownloadStatus(unlock *Unlock) {
	filename := manager.getCachePath(unlock)

	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		unlock.DownloadStatus = NotDownloaded
	} else if err != nil {
		unlock.DownloadStatus = NotDownloaded
		unlock.DownloadError = err
	} else {
		unlock.DownloadStatus = Downloaded
	}
}

func (manager *Manager) detectUnpackStatus(unlock *Unlock, user *UserData) {
	userUnlocks := settings.Get().UserUnlocks

	if !userUnlocks && user != unlock.Users[0] {
		user.UnpackStatus = unlock.Users[0].UnpackStatus
	}

	var profileName *string = nil
	if userUnlocks {
		profileName = &user.ProfileName
	}
	cookiePath := manager.getCookiePath(unlock, profileName)

	_, err := os.Stat(cookiePath)
	if os.IsNotExist(err) {
		user.UnpackStatus = NotUnpacked
	} else if err != nil {
		user.UnpackStatus = NotUnpacked
		user.UnpackError = err
	} else {
		user.UnpackStatus = Unpacked
	}
}

func (manager *Manager) getCachePath(unlock *Unlock) string {
	parts := strings.Split(unlock.DownloadUrl, "/")
	basename := parts[len(parts)-1]
	return filepath.Join(manager.DownloadDir, basename)
}

func (manager *Manager) getUnpackPath(unlock *Unlock, profileName *string) string {
	packName := fmt.Sprintf("%s Unlocks", unlock.RpgName)
	if profileName != nil {
		packName += fmt.Sprintf(" - %s", *profileName)
	}
	re := regexp.MustCompile(`[<>:"/\\|?*]`)
	packName = re.ReplaceAllLiteralString(packName, "_")

	return filepath.Join(settings.Get().SmDataDir, "Songs", packName)
}

func (manager *Manager) getCookiePath(unlock *Unlock, profileName *string) string {
	parts := strings.Split(unlock.DownloadUrl, "/")
	basename := parts[len(parts)-1]
	cookieName := basename + "-unpacked.txt"
	return filepath.Join(manager.getUnpackPath(unlock, profileName), cookieName)
}

func (manager *Manager) download(unlock *Unlock) {
	unlock.DownloadStatus = Downloading
	unlock.DownloadError = nil

	filename := manager.getCachePath(unlock)
	download := Fetch(unlock.DownloadUrl, filename)

	for info := range download.Progress {
		unlock.DownloadSize = info.TotalSize
		unlock.DownloadProgress = info.Downloaded
		if info.Error != nil {
			unlock.DownloadStatus = NotDownloaded
			unlock.DownloadError = info.Error
		}
		manager.updateCallback(unlock)
	}

	if unlock.DownloadError == nil {
		unlock.DownloadStatus = Downloaded
	}
	manager.updateCallback(unlock)
}

func (manager *Manager) unpack(unlock *Unlock) {
	for _, user := range unlock.Users {
		user.UnpackStatus = Unpacking
		user.UnpackError = nil
	}
	manager.updateCallback(unlock)

	filename := manager.getCachePath(unlock)
	unpackDir := manager.getUnpackPath(unlock, nil)

	err := unzip(filename, unpackDir)
	if err != nil {
		for _, user := range unlock.Users {
			user.UnpackStatus = NotUnpacked
			user.UnpackError = err
		}
		manager.updateCallback(unlock)
		return
	}

	cookiePath := manager.getCookiePath(unlock, nil)
	os.WriteFile(cookiePath, []byte(""), 0600)

	for _, user := range unlock.Users {
		user.UnpackStatus = Unpacked
	}
	manager.updateCallback(unlock)
}

func (manager *Manager) unpackUser(unlock *Unlock, user *UserData) {
	user.UnpackStatus = Unpacking
	user.UnpackError = nil
	manager.updateCallback(unlock)

	filename := manager.getCachePath(unlock)
	unpackDir := manager.getUnpackPath(unlock, &user.ProfileName)

	err := unzip(filename, unpackDir)
	if err != nil {
		user.UnpackStatus = NotUnpacked
		user.UnpackError = err
		manager.updateCallback(unlock)
		return
	}

	cookiePath := manager.getCookiePath(unlock, &user.ProfileName)
	os.WriteFile(cookiePath, []byte(""), 0600)

	user.UnpackStatus = Unpacked
	manager.updateCallback(unlock)
}

func (manager *Manager) refresh(unlock *Unlock) {
	for _, user := range unlock.Users {
		manager.detectUnpackStatus(unlock, user)
	}
	manager.updateCallback(unlock)
}
