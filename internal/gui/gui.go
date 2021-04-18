package gui

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
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

	logsMenuItem := fyne.NewMenuItem("StepMania Logs", nil)
	logsMenuItem.ChildMenu = fyne.NewMenu(
		"",
		fyne.NewMenuItem("info.txt", func() {
			app.openSMLog("info.txt")
		}),
		fyne.NewMenuItem("log.txt", func() {
			app.openSMLog("log.txt")
		}),
		fyne.NewMenuItem("timelog.txt", func() {
			app.openSMLog("timelog.txt")
		}),
		fyne.NewMenuItem("userlog.txt", func() {
			app.openSMLog("userlog.txt")
		}),
	)

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
			logsMenuItem,
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

	if settings.Get().FirstLaunch {
		app.showFirstLaunchDialog()
	}

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

func (app *App) openSMLog(filename string) {
	logPath := filepath.Join(settings.Get().SmDataDir, "Logs", filename)

	_, err := os.Stat(logPath)
	if err != nil {
		dialog.ShowError(err, app.mainWin)
		return
	}

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", logPath)
	} else {
		cmd = exec.Command("xdg-open", logPath)
	}

	err = cmd.Run()
	if err != nil {
		dialog.ShowError(err, app.mainWin)
	}
}
