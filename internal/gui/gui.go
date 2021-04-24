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
	app           fyne.App
	mainWin       fyne.Window
	unlockManager *unlocks.Manager
	unlockWidget  *UnlockWidget
	launchButton  *widget.Button
	session       *session.Session
	autolaunch    bool
}

func NewApp(unlockManager *unlocks.Manager, autolaunch bool) *App {
	app := &App{
		app:           app.New(),
		unlockManager: unlockManager,
		autolaunch:    autolaunch,
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
			filename := filepath.Join(settings.Get().SmDataDir, "Logs", "info.txt")
			app.viewLogfile(filename)
		}),
		fyne.NewMenuItem("log.txt", func() {
			filename := filepath.Join(settings.Get().SmDataDir, "Logs", "log.txt")
			app.viewLogfile(filename)
		}),
		fyne.NewMenuItem("timelog.txt", func() {
			filename := filepath.Join(settings.Get().SmDataDir, "Logs", "timelog.txt")
			app.viewLogfile(filename)
		}),
		fyne.NewMenuItem("userlog.txt", func() {
			filename := filepath.Join(settings.Get().SmDataDir, "Logs", "userlog.txt")
			app.viewLogfile(filename)
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
				go app.maybeQuit()
			}),
		),
		fyne.NewMenu(
			"View",
			logsMenuItem,
			fyne.NewMenuItem("Launcher Log", func() {
				cacheDir, err := os.UserCacheDir()
				if err != nil {
					dialog.ShowError(err, app.mainWin)
					return
				}

				filename := filepath.Join(cacheDir, "groovestats-launcher", "log.txt")
				app.viewLogfile(filename)
			}),
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

	app.launchButton = widget.NewButton("Launch StepMania", nil)
	app.launchButton.OnTapped = app.launchSM
	app.launchButton.Importance = widget.HighImportance

	app.unlockWidget = NewUnlockWidget(unlockManager)

	app.mainWin.SetContent(container.NewVScroll(container.NewVBox(
		app.unlockWidget.vbox,
		layout.NewSpacer(),
		container.NewPadded(app.launchButton),
	)))

	app.mainWin.CenterOnScreen()
	app.mainWin.Show()

	app.mainWin.SetCloseIntercept(func() { go app.maybeQuit() })

	if settings.Get().FirstLaunch && !autolaunch {
		app.showFirstLaunchDialog()
	}

	return app
}

func (app *App) Run() {
	if app.autolaunch {
		app.launchSM()
	}

	app.app.Run()
}

func (app *App) launchSM() {
	session, err := session.Launch(app.unlockManager)
	if err != nil {
		dialog.ShowError(err, app.mainWin)
		return
	}

	app.session = session
	app.launchButton.Disable()

	go func() {
		session.Wait()
		app.session = nil
		app.launchButton.Enable()

		if app.autolaunch && !app.unlockManager.HasPending() {
			app.mainWin.Close()
		}
	}()
}

func (app *App) maybeQuit() {
	session := app.session
	ch := make(chan bool, 10)

	if session != nil {
		confirmDialog := dialog.NewConfirm(
			"Stop StepMania?",
			"Closing the launcher will stop StepMania as well.",
			func(confirmed bool) {
				ch <- confirmed
			},
			app.mainWin,
		)
		confirmDialog.SetConfirmText("Stop StepMania")
		confirmDialog.SetDismissText("Keep Running")
		confirmDialog.Show()

		confirmed := <-ch
		if confirmed {
			session.Kill()
		} else {
			return
		}
	}

	if app.unlockManager.HasPending() {
		confirmDialog := dialog.NewConfirm(
			"Discard unlocks?",
			"Closing the launcher will also discard pending unlocks.\n"+
				"If you want to get them after closing the launcher"+
				" you will have to download them from the RPG website.",
			func(confirmed bool) {
				ch <- confirmed
			},
			app.mainWin,
		)
		confirmDialog.SetConfirmText("Discard Unlocks")
		confirmDialog.SetDismissText("Keep Running")
		confirmDialog.Show()

		confirmed := <-ch
		if !confirmed {
			return
		}
	}

	app.mainWin.Close()
}

func (app *App) showStatisticsDialog() {
	message := fmt.Sprintf("GET /new-session.php: %d\n", stats.GsNewSessionCount)
	message += fmt.Sprintf("GET /player-scores.php: %d\n", stats.GsPlayerScoresCount)
	message += fmt.Sprintf("GET /player-scores.php (cached): %d\n", stats.GsPlayerScoresCachedCount)
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

func (app *App) viewLogfile(filename string) {
	_, err := os.Stat(filename)
	if err != nil {
		dialog.ShowError(err, app.mainWin)
		return
	}

	var cmd *exec.Cmd

	// Open the file with the default application
	if runtime.GOOS == "windows" {
		cmd = exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", filename)
	} else {
		cmd = exec.Command("xdg-open", filename)
	}

	err = cmd.Run()
	if err != nil {
		dialog.ShowError(err, app.mainWin)
	}
}
