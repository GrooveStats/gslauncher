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
		fileDialog := dialog.NewFileOpen(func(file fyne.URIReadCloser, err error) {
			if err != nil || file == nil {
				return
			}

			path := filepath.FromSlash(file.URI().Path())
			data.SmExePath = path
			smExeButton.SetText(path)
		}, app.mainWin)
		uri, err := pathToUrl(filepath.Dir(data.SmExePath))
		if err == nil {
			fileDialog.SetLocation(uri)
		}
		fileDialog.Resize(fyne.NewSize(700, 500))
		fileDialog.Show()
	}
	smExeButton.SetText(abbreviatePath(data.SmExePath))

	smDirButton := widget.NewButton("Select", nil)
	smDirButton.OnTapped = func() {
		fileDialog := dialog.NewFolderOpen(func(dir fyne.ListableURI, err error) {
			if err != nil || dir == nil {
				return
			}

			path := filepath.FromSlash(dir.Path())

			_, err = os.Stat(filepath.Join(path, "Save"))
			if err != nil {
				err = errors.New("invalid data directory: Save folder missing")
				dialog.ShowError(err, app.mainWin)
				return
			}

			_, err = os.Stat(filepath.Join(path, "Songs"))
			if err != nil {
				err = errors.New("invalid data directory: Songs folder missing")
				dialog.ShowError(err, app.mainWin)
				return
			}

			data.SmDataDir = path
			smDirButton.SetText(path)
		}, app.mainWin)
		uri, err := pathToUrl(data.SmDataDir)
		if err == nil {
			fileDialog.SetLocation(uri)
		}
		fileDialog.Resize(fyne.NewSize(700, 500))
		fileDialog.Show()
	}
	smDirButton.SetText(abbreviatePath(data.SmDataDir))

	smDirFormItem := widget.NewFormItem("StepMania 5 Data Directory", smDirButton)
	smDirFormItem.HintText = "The folder containing Save and Songs"

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

	form := widget.NewForm(
		widget.NewFormItem("StepMania 5 Executable", smExeButton),
		smDirFormItem,
		autoDownloadFormItem,
		widget.NewFormItem("Separate Unlocks by User", userUnlocksCheck),
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

	firstLaunchDialog := dialog.NewCustom("Welcome!", "Save", container.NewVBox(
		welcomeMessage,
		widget.NewSeparator(),
		form,
	), app.mainWin)
	firstLaunchDialog.SetOnClosed(func() {
		settings.Update(data)

		err := settings.Save()
		if err != nil {
			dialog.ShowError(err, app.mainWin)
		}
	})
	firstLaunchDialog.Show()
}

func (app *App) showSettingsDialog() {
	data := settings.Get()

	form := app.getSettingsForm(&data)

	settingsDialog := dialog.NewCustomConfirm("Settings", "Save", "Cancel", form, func(save bool) {
		if save {
			settings.Update(data)

			err := settings.Save()
			if err != nil {
				dialog.ShowError(err, app.mainWin)
			}

			app.unlockWidget.Refresh()
		}
	}, app.mainWin)
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

	fakeGsCheck := widget.NewCheck("", func(checked bool) {
		data.FakeGs = checked

		if checked {
			fakeGsNetworkErrorCheck.Enable()
			fakeGsNewSessionResultSelect.Enable()
			fakeGsSubmitResultSelect.Enable()
			fakeGsRpgCheck.Enable()
			fakeGsNetDelayEntry.Enable()
		} else {
			fakeGsNetworkErrorCheck.Disable()
			fakeGsNewSessionResultSelect.Disable()
			fakeGsSubmitResultSelect.Disable()
			fakeGsRpgCheck.Disable()
			fakeGsNetDelayEntry.Disable()
		}
	})
	fakeGsCheck.SetChecked(data.FakeGs)

	if !data.FakeGs {
		fakeGsNetworkErrorCheck.Disable()
		fakeGsNewSessionResultSelect.Disable()
		fakeGsSubmitResultSelect.Disable()
		fakeGsRpgCheck.Disable()
		fakeGsNetDelayEntry.Disable()
	}

	gsUrlEntry := widget.NewEntry()
	gsUrlEntry.Text = data.GrooveStatsUrl
	gsUrlEntry.OnChanged = func(url string) {
		data.GrooveStatsUrl = url
	}

	items := []*widget.FormItem{
		widget.NewFormItem("Simulate GrooveStats Requests", fakeGsCheck),
		widget.NewFormItem(">> Network Error", fakeGsNetworkErrorCheck),
		widget.NewFormItem(">> Network Delay (Seconds)", fakeGsNetDelayEntry),
		widget.NewFormItem(">> New Session Result", fakeGsNewSessionResultSelect),
		widget.NewFormItem(">> Score Submit Result", fakeGsSubmitResultSelect),
		widget.NewFormItem(">> RPG active", fakeGsRpgCheck),
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
