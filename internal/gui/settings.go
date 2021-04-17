package gui

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/GrooveStats/gslauncher/internal/settings"
)

func (app *App) getSettingsFormItems(data *settings.Settings) []*widget.FormItem {
	smExeButton := widget.NewButton("Select", nil)
	smExeButton.OnTapped = func() {
		fileDialog := dialog.NewFileOpen(func(file fyne.URIReadCloser, err error) {
			if err != nil || file == nil {
				return
			}

			path := file.URI().Path()
			data.SmExePath = path
			smExeButton.SetText(path)
		}, app.mainWin)
		fileDialog.Resize(fyne.NewSize(700, 500))
		fileDialog.Show()
	}
	smExeButton.SetText(data.SmExePath)

	smDirButton := widget.NewButton("Select", nil)
	smDirButton.OnTapped = func() {
		fileDialog := dialog.NewFolderOpen(func(dir fyne.ListableURI, err error) {
			if err != nil || dir == nil {
				return
			}

			path := dir.Path()
			data.SmDataDir = path
			smDirButton.SetText(path)
		}, app.mainWin)
		fileDialog.Resize(fyne.NewSize(700, 500))
		fileDialog.Show()
	}
	smDirButton.SetText(data.SmDataDir)

	autoDownloadCheck := widget.NewCheck("", func(checked bool) {
		data.AutoDownload = checked
	})
	autoDownloadCheck.SetChecked(data.AutoDownload)

	autoUnpackCheck := widget.NewCheck("", func(checked bool) {
		data.AutoUnpack = checked
	})
	autoUnpackCheck.SetChecked(data.AutoUnpack)

	userUnlocksCheck := widget.NewCheck("", func(checked bool) {
		data.UserUnlocks = checked
	})
	userUnlocksCheck.SetChecked(data.UserUnlocks)

	return []*widget.FormItem{
		widget.NewFormItem("StepMania 5 Executable", smExeButton),
		widget.NewFormItem("StepMania 5 Data Directory", smDirButton),
		widget.NewFormItem("Auto-Download Unlocks", autoDownloadCheck),
		widget.NewFormItem("Auto-Unpack Unlocks", autoUnpackCheck),
		widget.NewFormItem("Separate Unlocks by User", userUnlocksCheck),
	}
}

func (app *App) showFirstLaunchDialog() {
	data := settings.Get()

	message := "Thank you for using the GrooveStat Launcher!\n"
	message += "Please review the settings below before you continue."
	welcomeMessage := widget.NewLabel(message)
	welcomeMessage.Wrapping = fyne.TextWrapWord
	welcomeMessage.Alignment = fyne.TextAlignCenter

	form := widget.NewForm()
	form.Items = app.getSettingsFormItems(&data)

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

	items := app.getSettingsFormItems(&data)

	if data.Debug {
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

		items = append(
			items,
			widget.NewFormItem("[DEBUG] Simulate GrooveStats Requests", fakeGsCheck),
			widget.NewFormItem("[DEBUG] >> Network Error", fakeGsNetworkErrorCheck),
			widget.NewFormItem("[DEBUG] >> Network Delay (Seconds)", fakeGsNetDelayEntry),
			widget.NewFormItem("[DEBUG] >> New Session Result", fakeGsNewSessionResultSelect),
			widget.NewFormItem("[DEBUG] >> Score Submit Result", fakeGsSubmitResultSelect),
			widget.NewFormItem("[DEBUG] >> RPG active", fakeGsRpgCheck),
		)

		gsUrlEntry := widget.NewEntry()
		gsUrlEntry.Text = data.GrooveStatsUrl
		gsUrlEntry.OnChanged = func(url string) {
			data.GrooveStatsUrl = url
		}
		items = append(items, widget.NewFormItem("[DEBUG] GrooveStats URL", gsUrlEntry))
	}

	dialog.ShowForm("Settings", "Save", "Cancel", items, func(save bool) {
		if save {
			settings.Update(data)

			err := settings.Save()
			if err != nil {
				dialog.ShowError(err, app.mainWin)
			}
		}
	}, app.mainWin)
}