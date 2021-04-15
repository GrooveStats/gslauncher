package gui

import (
	"fmt"
	"net/url"
	"runtime"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/GrooveStats/gslauncher/internal/session"
	"github.com/GrooveStats/gslauncher/internal/settings"
	"github.com/GrooveStats/gslauncher/internal/stats"
	"github.com/GrooveStats/gslauncher/internal/unlocks"
	"github.com/GrooveStats/gslauncher/internal/version"
)

type App struct {
	app          fyne.App
	mainWin      fyne.Window
	unlockWidget *UnlockWidget
	session      *session.Session
}

func NewApp(unlockManager *unlocks.Manager) *App {
	app := &App{
		app: app.New(),
	}

	app.app.Settings().SetTheme(theme.DarkTheme())
	app.app.SetIcon(groovestatsLogo)

	appName := "GrooveStats Launcher"
	if settings.Get().Debug {
		appName += " (debug)"
	}

	app.mainWin = app.app.NewWindow(appName)
	app.mainWin.Resize(fyne.NewSize(800, 600))

	app.mainWin.SetMainMenu(fyne.NewMainMenu(
		fyne.NewMenu(
			"File",
			fyne.NewMenuItem("Settings", func() {
				app.showSettingsDialog()
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Quit", func() {
				app.maybeQuit()
			}),
		),
		fyne.NewMenu(
			"View",
			fyne.NewMenuItem("Statistics", func() {
				app.showStatisticsDialog()
			}),
		),
		fyne.NewMenu(
			"Help",
			fyne.NewMenuItem("Setup", func() {
				url, err := url.Parse("https://github.com/GrooveStats/gslauncher#readme")
				if err != nil {
					return
				}

				app.app.OpenURL(url)
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("About", func() {
				app.showAboutDialog()
			}),
		),
	))

	launchButton := widget.NewButton("Launch StepMania", nil)
	launchButton.OnTapped = func() {
		session, err := session.Launch(unlockManager)
		if err != nil {
			dialog.ShowError(err, app.mainWin)
			return
		}

		app.session = session
		launchButton.Disable()

		go func() {
			session.Wait()
			app.session = nil
			launchButton.Enable()
		}()
	}
	launchButton.Importance = widget.HighImportance

	app.unlockWidget = NewUnlockWidget(unlockManager)

	app.mainWin.SetContent(container.NewVBox(
		app.unlockWidget.vbox,
		layout.NewSpacer(),
		container.NewPadded(launchButton),
	))

	app.mainWin.CenterOnScreen()
	app.mainWin.Show()

	app.mainWin.SetCloseIntercept(app.maybeQuit)

	return app
}

func (app *App) Run() {
	app.app.Run()
}

func (app *App) maybeQuit() {
	session := app.session

	if session != nil {
		confirmDialog := dialog.NewConfirm(
			"Stop StepMania?",
			"Closing the launcher will stop StepMania as well.",
			func(confirmed bool) {
				if confirmed {
					app.session.Kill()
					app.mainWin.Close()
				}
			},
			app.mainWin,
		)
		confirmDialog.SetConfirmText("Stop StepMania")
		confirmDialog.SetDismissText("Keep Running")
		confirmDialog.Show()
	} else {
		app.mainWin.Close()
	}
}

func (app *App) showSettingsDialog() {
	data := settings.Get()

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

func (app *App) showStatisticsDialog() {
	message := fmt.Sprintf("GET /new-session.php: %d\n", stats.GsNewSessionCount)
	message += fmt.Sprintf("GET /player-scores.php: %d\n", stats.GsPlayerScoresCount)
	message += fmt.Sprintf("GET /player-leaderboards.php: %d\n", stats.GsPlayerLeaderboardsCount)
	message += fmt.Sprintf("POST /score-submit.php: %d\n", stats.GsScoreSubmitCount)

	dialog.ShowInformation("Statistics", message, app.mainWin)
}

func (app *App) showAboutDialog() {
	message := fmt.Sprintf(
		"GrooveStats Launcher\n%s (%s %s)",
		version.Formatted(), runtime.GOOS, runtime.GOARCH,
	)
	if settings.Get().Debug {
		message += "\ndebug"
	}

	dialog.ShowInformation("About", message, app.mainWin)
}
