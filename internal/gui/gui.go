package gui

import (
	"os/exec"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/archiveflax/gslauncher/internal/settings"
)

type App struct {
	app     fyne.App
	mainWin fyne.Window
}

func NewApp() *App {
	app := &App{
		app: app.New(),
	}

	app.app.Settings().SetTheme(theme.DarkTheme())

	app.mainWin = app.app.NewWindow("GrooveStats Launcher")
	app.mainWin.Resize(fyne.NewSize(800, 600))

	app.mainWin.SetMainMenu(fyne.NewMainMenu(
		fyne.NewMenu("File", fyne.NewMenuItem("Settings", func() {
			app.showSettingsDialog()
		})),
	))

	app.mainWin.SetContent(container.NewVBox(
		widget.NewButton("Launch StepMania", func() {
			cmd := exec.Command(settings.Get().SmExePath)

			err := cmd.Start()
			if err != nil {
				dialog.ShowError(err, app.mainWin)
				return
			}

			go cmd.Wait()
		}),
	))

	app.mainWin.CenterOnScreen()
	app.mainWin.Show()

	return app
}

func (app *App) Run() {
	app.app.Run()
}

func (app *App) showSettingsDialog() {
	data := settings.Get()

	smExeButton := widget.NewButton("Select", nil)
	smExeButton.OnTapped = func() {
		dialog.ShowFileOpen(func(file fyne.URIReadCloser, err error) {
			if err != nil || file == nil {
				return
			}

			path := file.URI().Path()
			data.SmExePath = path
			smExeButton.SetText(path)
		}, app.mainWin)
	}
	smExeButton.SetText(data.SmExePath)

	smDirButton := widget.NewButton("Select", nil)
	smDirButton.OnTapped = func() {
		dialog.ShowFolderOpen(func(dir fyne.ListableURI, err error) {
			if err != nil || dir == nil {
				return
			}

			path := dir.Path()
			data.SmDataDir = path
			smDirButton.SetText(path)
		}, app.mainWin)
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

	items := []*widget.FormItem{
		widget.NewFormItem("StepMania 5 Executable", smExeButton),
		widget.NewFormItem("StepMania 5 Data Directory", smDirButton),
		widget.NewFormItem("Auto-Download Unlocks", autoDownloadCheck),
		widget.NewFormItem("Auto-Unpack Unlocks", autoUnpackCheck),
		widget.NewFormItem("Separate Unlocks by User", userUnlocksCheck),
	}

	if data.Debug {
		fakeGsNetworkErrorCheck := widget.NewCheck("", func(checked bool) {
			data.FakeGsNetworkError = checked
		})
		fakeGsNetworkErrorCheck.SetChecked(data.FakeGsNetworkError)

		fakeGsDdosCheck := widget.NewCheck("", func(checked bool) {
			data.FakeGsDdos = checked
		})
		fakeGsDdosCheck.SetChecked(data.FakeGsDdos)

		options := []string{"score-added", "improved", "score-not-improved", "score-not-ranked"}
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
				fakeGsDdosCheck.Enable()
				fakeGsSubmitResultSelect.Enable()
				fakeGsRpgCheck.Enable()
			} else {
				fakeGsNetworkErrorCheck.Disable()
				fakeGsDdosCheck.Disable()
				fakeGsSubmitResultSelect.Disable()
				fakeGsRpgCheck.Disable()
			}
		})
		fakeGsCheck.SetChecked(data.FakeGs)

		if !data.FakeGs {
			fakeGsNetworkErrorCheck.Disable()
			fakeGsDdosCheck.Disable()
			fakeGsSubmitResultSelect.Disable()
			fakeGsRpgCheck.Disable()
		}

		items = append(
			items,
			widget.NewFormItem("[DEBUG] Simulate GrooveStats Requests", fakeGsCheck),
			widget.NewFormItem("[DEBUG] >> Network Error", fakeGsNetworkErrorCheck),
			widget.NewFormItem("[DEBUG] >> DDoS Protection", fakeGsDdosCheck),
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
