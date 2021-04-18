package unlocks

import (
	"github.com/GrooveStats/gslauncher/internal/settings"
)

type actionDownload struct{}
type actionRefresh struct{}
type actionUnpack struct{ user *UserData }

func (manager *Manager) processQueue(unlock *Unlock) {
	for action := range unlock.queue {
		switch a := action.(type) {
		case actionDownload:
			manager.doDownload(unlock)
		case actionRefresh:
			manager.refresh(unlock)
		case actionUnpack:
			manager.doUnpack(unlock, a.user)
		}
	}
}

func (manager *Manager) doDownload(unlock *Unlock) {
	if unlock.DownloadStatus != NotDownloaded {
		return
	}

	manager.download(unlock)
}

func (manager *Manager) doUnpack(unlock *Unlock, user *UserData) {
	if unlock.DownloadStatus != Downloaded {
		return
	}

	if settings.Get().UserUnlocks {
		if user.UnpackStatus != NotUnpacked {
			return
		}

		if user == nil {
			return
		}

		manager.unpackUser(unlock, user)
	} else {
		if unlock.Users[0].UnpackStatus != NotUnpacked {
			return
		}

		manager.unpack(unlock)
	}
}

func (unlock *Unlock) QueueDownload() {
	unlock.queue <- actionDownload{}
}

func (unlock *Unlock) QueueRefresh() {
	unlock.queue <- actionRefresh{}
}

func (unlock *Unlock) QueueUnpack(user *UserData) {
	unlock.queue <- actionUnpack{user: user}
}
