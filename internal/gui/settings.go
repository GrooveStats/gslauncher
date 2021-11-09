package gui

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"github.com/GrooveStats/gslauncher/internal/settings"
)

func pathToUrl(fpath string) (fyne.ListableURI, error) {
	fpath = filepath.ToSlash(fpath)

	uri, err := storage.ParseURI("file://" + fpath)
	if err != nil {
		return nil, err
	}

	return storage.ListerForURI(uri)
}

func abbreviatePath(fpath string) string {
	if runtime.GOOS == "windows" {
		configDir, err := os.UserConfigDir()
		if err == nil {
			if strings.HasPrefix(fpath, configDir+string(os.PathSeparator)) {
				rel, err := filepath.Rel(configDir, fpath)
				if err == nil {
					fpath = "%AppData%\\" + rel
				}
			}
		}
	} else {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			if strings.HasPrefix(fpath, homeDir+string(os.PathSeparator)) {
				rel, err := filepath.Rel(homeDir, fpath)
				if err == nil {
					fpath = "~/" + rel
				}
			}
		}
	}

	if len(fpath) > 55 {
		fpath = fpath[:52] + "..."
	}

	return fpath
}

func (app *App) getSettingsForm(data *settings.Settings) fyne.CanvasObject {
	smExeButton := widget.NewButton("Select", nil)
	smExeButton.OnTapped = func() {
		if runtime.GOOS == "darwin" {
			fileDialog := dialog.NewFolderOpen(func(dir fyne.ListableURI, err error) {
				if err != nil || dir == nil {
					return
				}

				path := filepath.FromSlash(dir.Path())

				if !strings.HasSuffix(path, ".app") {
					err = errors.New("invalid application: must be a .app directory")
					dialog.ShowError(err, app.mainWin)
					return
				}

				data.SmExePath = path
				smExeButton.SetText(abbreviatePath(path))
			}, app.mainWin)
			uri, err := pathToUrl(filepath.Dir(data.SmExePath))
			if err == nil {
				fileDialog.SetLocation(uri)
			}
			fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".app"}))
			fileDialog.Resize(fyne.NewSize(700, 500))
			fileDialog.Show()
		} else {
			fileDialog := dialog.NewFileOpen(func(file fyne.URIReadCloser, err error) {
				if err != nil || file == nil {
					return
				}

				path := filepath.FromSlash(file.URI().Path())
				data.SmExePath = path
				smExeButton.SetText(abbreviatePath(path))
			}, app.mainWin)
			uri, err := pathToUrl(filepath.Dir(data.SmExePath))
			if err == nil {
				fileDialog.SetLocation(uri)
			}
			fileDialog.Resize(fyne.NewSize(700, 500))
			fileDialog.Show()
		}
	}
	smExeButton.SetText(abbreviatePath(data.SmExePath))

	smExeButtonFormItem := widget.NewFormItem("StepMania 5 Executable", smExeButton)
	if runtime.GOOS == "darwin" {
		smExeButtonFormItem.Text = "StepMania 5 App"
	}
	smExeButtonFormItem.HintText = "Currently supported are SM 5.0, 5.1 and 5.3 (Outfox)"

	smSaveDirButton := widget.NewButton("Select", nil)
	smSaveDirButton.OnTapped = func() {
		fileDialog := dialog.NewFolderOpen(func(dir fyne.ListableURI, err error) {
			if err != nil || dir == nil {
				return
			}

			path := filepath.FromSlash(dir.Path())

			_, err = os.Stat(filepath.Join(path, "Preferences.ini"))
			if err != nil {
				err = errors.New("invalid save directory: Preferences.ini not found")
				dialog.ShowError(err, app.mainWin)
				return
			}

			data.SmSaveDir = path
			smSaveDirButton.SetText(abbreviatePath(path))
		}, app.mainWin)
		uri, err := pathToUrl(data.SmSaveDir)
		if err == nil {
			fileDialog.SetLocation(uri)
		}
		fileDialog.Resize(fyne.NewSize(700, 500))
		fileDialog.Show()
	}
	smSaveDirButton.SetText(abbreviatePath(data.SmSaveDir))

	smSaveDirFormItem := widget.NewFormItem("StepMania 5 Save Directory", smSaveDirButton)
	smSaveDirFormItem.HintText = "The folder containing StepMania's Preferences.ini"

	smLogsDirButton := widget.NewButton("Select", nil)
	smLogsDirButton.OnTapped = func() {
		fileDialog := dialog.NewFolderOpen(func(dir fyne.ListableURI, err error) {
			if err != nil || dir == nil {
				return
			}

			path := filepath.FromSlash(dir.Path())

			_, err = os.Stat(filepath.Join(path, "info.txt"))
			if err != nil {
				err = errors.New("invalid logs directory: info.txt not found")
				dialog.ShowError(err, app.mainWin)
				return
			}

			data.SmLogsDir = path
			smLogsDirButton.SetText(abbreviatePath(path))
		}, app.mainWin)
		uri, err := pathToUrl(data.SmLogsDir)
		if err == nil {
			fileDialog.SetLocation(uri)
		}
		fileDialog.Resize(fyne.NewSize(700, 500))
		fileDialog.Show()
	}
	smLogsDirButton.SetText(abbreviatePath(data.SmLogsDir))

	smLogsDirFormItem := widget.NewFormItem("StepMania 5 Logs Directory", smLogsDirButton)

	smSongsDirButton := widget.NewButton("Select", nil)
	smSongsDirButton.OnTapped = func() {
		fileDialog := dialog.NewFolderOpen(func(dir fyne.ListableURI, err error) {
			if err != nil || dir == nil {
				return
			}

			path := filepath.FromSlash(dir.Path())

			data.SmSongsDir = path
			smSongsDirButton.SetText(abbreviatePath(path))
		}, app.mainWin)
		uri, err := pathToUrl(data.SmSongsDir)
		if err == nil {
			fileDialog.SetLocation(uri)
		}
		fileDialog.Resize(fyne.NewSize(700, 500))
		fileDialog.Show()
	}
	smSongsDirButton.SetText(abbreviatePath(data.SmSongsDir))

	smSongsDirFormItem := widget.NewFormItem("StepMania 5 Songs Directory", smSongsDirButton)
	smSongsDirFormItem.HintText = "Unlocked RPG songs will be stored here"

	options := []string{"Off", "Download Only", "Download and Unpack"}
	autoDownloadSelect := widget.NewSelect(options, func(selected string) {
		switch selected {
		case "Off":
			data.AutoDownloadMode = settings.AutoDownloadOff
		case "Download Only":
			data.AutoDownloadMode = settings.AutoDownloadOnly
		case "Download and Unpack":
			data.AutoDownloadMode = settings.AutoDownloadAndUnpack
		}
	})
	switch data.AutoDownloadMode {
	case settings.AutoDownloadOff:
		autoDownloadSelect.SetSelected("Off")
	case settings.AutoDownloadOnly:
		autoDownloadSelect.SetSelected("Download Only")
	case settings.AutoDownloadAndUnpack:
		autoDownloadSelect.SetSelected("Download and Unpack")
	}

	autoDownloadFormItem := widget.NewFormItem("Automatically Download\nUnlocked RPG Songs", autoDownloadSelect)
	autoDownloadFormItem.HintText = "Can negatively impact game performance when enabled!"

	userUnlocksCheck := widget.NewCheck("", func(checked bool) {
		data.UserUnlocks = checked
	})
	userUnlocksCheck.SetChecked(data.UserUnlocks)

	autoLaunchCheck := widget.NewCheck("", func(checked bool) {
		data.AutoLaunch = checked
	})
	autoLaunchCheck.SetChecked(data.AutoLaunch)

	gsUrlEntry := widget.NewEntry()
	gsUrlEntry.Text = data.GrooveStatsUrl
	gsUrlEntry.OnChanged = func(url string) {
		data.GrooveStatsUrl = url
	}

	gsUrlFormItem := widget.NewFormItem("GrooveStats API endpoint", gsUrlEntry)
	gsUrlFormItem.HintText = "Do not modify unless you know what you are doing!"

	form := widget.NewForm(
		smExeButtonFormItem,
		smSaveDirFormItem,
		smLogsDirFormItem,
		smSongsDirFormItem,
		autoDownloadFormItem,
		widget.NewFormItem("Separate Unlocks by User", userUnlocksCheck),
		widget.NewFormItem("Launch StepMania at Startup", autoLaunchCheck),
		gsUrlFormItem,
	)

	return form
}

func (app *App) showFirstLaunchDialog() {
	data := settings.Get()

	message := "Thank you for using the GrooveStat Launcher!\n"
	message += "Please review the settings below before you continue."
	welcomeMessage := widget.NewLabel(message)
	welcomeMessage.Wrapping = fyne.TextWrapWord
	welcomeMessage.Alignment = fyne.TextAlignCenter

	form := app.getSettingsForm(&data)
	content := container.NewVScroll(
		container.NewVBox(
			welcomeMessage,
			widget.NewSeparator(),
			form,
		),
	)

	firstLaunchDialog := dialog.NewCustom("Welcome!", "Save", content, app.mainWin)
	firstLaunchDialog.SetOnClosed(func() {
		settings.Update(data)

		err := settings.Save()
		if err != nil {
			dialog.ShowError(err, app.mainWin)
		}
	})
	firstLaunchDialog.Resize(fyne.NewSize(700, 560))
	firstLaunchDialog.Show()
}

func (app *App) showSettingsDialog() {
	data := settings.Get()

	form := app.getSettingsForm(&data)
	content := container.NewVScroll(form)

	settingsDialog := dialog.NewCustomConfirm("Settings", "Save", "Cancel", content, func(save bool) {
		if save {
			settings.Update(data)

			err := settings.Save()
			if err != nil {
				dialog.ShowError(err, app.mainWin)
			}

			app.unlockWidget.Refresh()
		}
	}, app.mainWin)
	settingsDialog.Resize(fyne.NewSize(700, 550))
	settingsDialog.Show()
}

func (app *App) showDebugSettingsDialog() {
	data := settings.Get()

	fakeGsNetworkErrorCheck := widget.NewCheck("", func(checked bool) {
		data.FakeGsNetworkError = checked
	})
	fakeGsNetworkErrorCheck.SetChecked(data.FakeGsNetworkError)

	fakeGsNetDelayEntry := widget.NewEntry()
	fakeGsNetDelayEntry.Validator = validation.NewRegexp(`^\d+$`, "Must contain a number")
	fakeGsNetDelayEntry.Text = strconv.Itoa(data.FakeGsNetworkDelay)
	fakeGsNetDelayEntry.OnChanged = func(s string) {
		n, err := strconv.Atoi(s)
		if err == nil {
			data.FakeGsNetworkDelay = n
		}
	}

	options := []string{"OK", "UNSUPPORTED_CHART_HASH", "MAINTENANCE"}
	fakeGsNewSessionResultSelect := widget.NewSelect(options, func(selected string) {
		data.FakeGsNewSessionResult = selected
	})
	fakeGsNewSessionResultSelect.SetSelected(data.FakeGsNewSessionResult)

	options = []string{"score-added", "improved", "score-not-improved", "chart-not-ranked"}
	fakeGsSubmitResultSelect := widget.NewSelect(options, func(selected string) {
		data.FakeGsSubmitResult = selected
	})
	fakeGsSubmitResultSelect.SetSelected(data.FakeGsSubmitResult)

	fakeGsRpgCheck := widget.NewCheck("", func(checked bool) {
		data.FakeGsRpg = checked
	})
	fakeGsRpgCheck.SetChecked(data.FakeGsRpg)

	fakeGsItlCheck := widget.NewCheck("", func(checked bool) {
		data.FakeGsItl = checked
	})
	fakeGsItlCheck.SetChecked(data.FakeGsItl)

	fakeGsCheck := widget.NewCheck("", func(checked bool) {
		data.FakeGs = checked

		if checked {
			fakeGsNetworkErrorCheck.Enable()
			fakeGsNewSessionResultSelect.Enable()
			fakeGsSubmitResultSelect.Enable()
			fakeGsRpgCheck.Enable()
			fakeGsItlCheck.Enable()
			fakeGsNetDelayEntry.Enable()
		} else {
			fakeGsNetworkErrorCheck.Disable()
			fakeGsNewSessionResultSelect.Disable()
			fakeGsSubmitResultSelect.Disable()
			fakeGsRpgCheck.Disable()
			fakeGsItlCheck.Disable()
			fakeGsNetDelayEntry.Disable()
		}
	})
	fakeGsCheck.SetChecked(data.FakeGs)

	if !data.FakeGs {
		fakeGsNetworkErrorCheck.Disable()
		fakeGsNewSessionResultSelect.Disable()
		fakeGsSubmitResultSelect.Disable()
		fakeGsRpgCheck.Disable()
		fakeGsItlCheck.Disable()
		fakeGsNetDelayEntry.Disable()
	}

	items := []*widget.FormItem{
		widget.NewFormItem("Simulate GrooveStats Requests", fakeGsCheck),
		widget.NewFormItem(">> Network Error", fakeGsNetworkErrorCheck),
		widget.NewFormItem(">> Network Delay (Seconds)", fakeGsNetDelayEntry),
		widget.NewFormItem(">> New Session Result", fakeGsNewSessionResultSelect),
		widget.NewFormItem(">> Score Submit Result", fakeGsSubmitResultSelect),
		widget.NewFormItem(">> RPG active", fakeGsRpgCheck),
		widget.NewFormItem(">> ITL active", fakeGsItlCheck),
		widget.NewFormItem("GrooveStats URL", gsUrlEntry),
	}

	formDialog := dialog.NewForm("Debug Settings", "Save", "Cancel", items, func(save bool) {
		if save {
			settings.Update(data)
		}
	}, app.mainWin)
	formDialog.Show()
	formDialog.Resize(fyne.NewSize(600, 300))
}
