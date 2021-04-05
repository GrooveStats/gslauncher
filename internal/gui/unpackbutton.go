package gui

import (
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/GrooveStats/gslauncher/internal/settings"
	"github.com/GrooveStats/gslauncher/internal/unlocks"
)

type unpackButton struct {
	widget.Button

	unlockManager *unlocks.Manager
	unlock        *unlocks.Unlock
}

func (button *unpackButton) Tapped(e *fyne.PointEvent) {
	if settings.Get().UserUnlocks {
		items := make([]*fyne.MenuItem, 0)

		for _, u := range button.unlock.Users {
			user := u

			switch user.UnpackStatus {
			case unlocks.NotUnpacked, unlocks.UnpackFailed:
				menuItem := fyne.NewMenuItem(
					"Unpack for "+user.ProfileName,
					func() {
						go button.unlockManager.UnpackUser(button.unlock, user)
					},
				)
				items = append(items, menuItem)
			}
		}

		sort.Slice(items, func(i, j int) bool {
			return items[i].Label < items[j].Label
		})

		menu := fyne.NewMenu("")
		menu.Items = items

		widget.ShowPopUpMenuAtPosition(
			menu,
			fyne.CurrentApp().Driver().CanvasForObject(button),
			e.AbsolutePosition,
		)
	} else {
		go button.unlockManager.Unpack(button.unlock)
	}
}

func newUnpackButton(unlockManager *unlocks.Manager, unlock *unlocks.Unlock) *unpackButton {
	button := &unpackButton{
		Button: widget.Button{
			Text: "Unpack",
			Icon: theme.FolderOpenIcon(),
		},

		unlockManager: unlockManager,
		unlock:        unlock,
	}

	button.ExtendBaseWidget(button)
	return button
}
