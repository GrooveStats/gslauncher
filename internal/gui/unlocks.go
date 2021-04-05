package gui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/GrooveStats/gslauncher/internal/unlocks"
)

type unlockInfo struct {
	unlock           *unlocks.Unlock
	downloadButton   *widget.Button
	downloadProgress *widget.ProgressBar
	unpackButton     *unpackButton
	unpackProgress   *widget.ProgressBarInfinite
	successIcon      *widget.Icon
	errorIcon        *widget.Icon
	errorLabel       *widget.Label
}

type UnlockWidget struct {
	unlockManager *unlocks.Manager
	unlockInfos   map[*unlocks.Unlock]*unlockInfo
	vbox          *fyne.Container
	emptyLabel    *widget.Label
}

func NewUnlockWidget(unlockManager *unlocks.Manager) *UnlockWidget {
	emptyLabel := widget.NewLabel("No unlocks yet.")
	emptyLabel.TextStyle = fyne.TextStyle{Italic: true}
	emptyLabel.Alignment = fyne.TextAlignCenter

	unlockWidget := &UnlockWidget{
		unlockManager: unlockManager,
		unlockInfos:   make(map[*unlocks.Unlock]*unlockInfo),
		vbox:          container.NewVBox(emptyLabel),
		emptyLabel:    emptyLabel,
	}

	unlockManager.SetUpdateCallback(unlockWidget.handleUpdate)

	return unlockWidget
}

func (unlockWidget *UnlockWidget) handleUpdate(unlock *unlocks.Unlock) {
	info, ok := unlockWidget.unlockInfos[unlock]

	if !ok {
		unlockWidget.emptyLabel.Hide()

		downloadButton := widget.NewButton("Download", func() {
			go unlockWidget.unlockManager.Download(unlock)
		})
		downloadButton.SetIcon(theme.DownloadIcon())

		unpackButton := newUnpackButton(unlockWidget.unlockManager, unlock)

		downloadProgress := widget.NewProgressBar()
		downloadProgress.Min = 0
		downloadProgress.Max = 1
		downloadProgress.SetValue(0)
		downloadProgress.TextFormatter = func() string {
			if unlock.DownloadSize == -1 {
				return "Connecting..."
			}

			return fmt.Sprintf(
				"%s / %s",
				formatBytes(unlock.DownloadProgress),
				formatBytes(unlock.DownloadSize),
			)
		}

		unpackProgress := widget.NewProgressBarInfinite()

		successIcon := widget.NewIcon(theme.ConfirmIcon())
		errorIcon := widget.NewIcon(theme.NewErrorThemedResource(theme.ErrorIcon()))

		errorLabel := widget.NewLabel("")
		errorLabel.Wrapping = fyne.TextWrapWord
		errorLabel.Alignment = fyne.TextAlignCenter

		vbox := container.NewVBox(
			container.NewHBox(
				widget.NewLabel(unlock.Description),
				layout.NewSpacer(),
				successIcon,
				errorIcon,
				downloadButton,
				unpackButton,
			),
			downloadProgress,
			unpackProgress,
			errorLabel,
		)
		unlockWidget.vbox.Add(vbox)
		unlockWidget.vbox.Add(widget.NewSeparator())

		info = &unlockInfo{
			unlock:           unlock,
			downloadButton:   downloadButton,
			downloadProgress: downloadProgress,
			unpackButton:     unpackButton,
			unpackProgress:   unpackProgress,
			successIcon:      successIcon,
			errorIcon:        errorIcon,
			errorLabel:       errorLabel,
		}
		unlockWidget.unlockInfos[unlock] = info
	}

	switch unlock.DownloadStatus {
	case unlocks.NotDownloaded:
		info.downloadButton.Show()
		info.downloadProgress.Hide()
		info.unpackButton.Hide()
		info.unpackProgress.Hide()
		info.unpackProgress.Stop()
		info.successIcon.Hide()
		info.errorIcon.Hide()
		info.errorLabel.Hide()
	case unlocks.Downloading:
		progress := float64(unlock.DownloadProgress) / float64(unlock.DownloadSize)

		info.downloadButton.Hide()
		info.downloadProgress.Show()
		info.downloadProgress.SetValue(progress)
		info.unpackButton.Hide()
		info.unpackProgress.Hide()
		info.unpackProgress.Stop()
		info.successIcon.Hide()
		info.errorIcon.Hide()
		info.errorLabel.Hide()
	case unlocks.Downloaded:
		info.downloadButton.Hide()
		info.downloadProgress.Hide()

		unpacked := true
		unpacking := false
		unpackErrors := make([]string, 0)

		for _, user := range unlock.Users {
			switch user.UnpackStatus {
			case unlocks.NotUnpacked:
				unpacked = false
			case unlocks.Unpacking:
				unpacking = true
			case unlocks.UnpackFailed:
				unpacked = false

				unpackErrors = append(
					unpackErrors,
					fmt.Sprintf("Unpack failed: %v", user.UnpackError),
				)
			}
		}

		if unpacked {
			info.unpackButton.Hide()
			info.successIcon.Show()
		} else {
			info.unpackButton.Show()
			info.successIcon.Hide()
		}

		if unpacking {
			info.unpackProgress.Show()
			info.unpackProgress.Start()
		} else {
			info.unpackProgress.Hide()
			info.unpackProgress.Stop()
		}

		if len(unpackErrors) > 0 {
			info.errorIcon.Show()
			info.errorLabel.Show()
			info.errorLabel.SetText(strings.Join(unpackErrors, "\n"))
		} else {
			info.errorIcon.Hide()
			info.errorLabel.Hide()
		}
	case unlocks.DownloadFailed:
		info.downloadButton.Show()
		info.downloadProgress.Hide()
		info.unpackButton.Hide()
		info.unpackProgress.Hide()
		info.unpackProgress.Stop()
		info.errorIcon.Show()
		info.errorLabel.Show()
		info.errorLabel.SetText(fmt.Sprintf("Download failed: %v", unlock.DownloadError))
	}
}

func formatBytes(n int64) string {
	switch {
	case n < 1024:
		return fmt.Sprintf("%d B", n)
	case n < 1024*1024:
		return fmt.Sprintf("%.1f KiB", float64(n)/1024)
	case n < 1024*1024*1024:
		return fmt.Sprintf("%.1f MiB", float64(n)/1024/1024)
	default:
		return fmt.Sprintf("%.1f GiB", float64(n)/1024/1024/1024)
	}
}
