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

	unlock *unlocks.Unlock
}

func (button *unpackButton) Tapped(e *fyne.PointEvent) {
	if settings.Get().UserUnlocks {
		items := make([]*fyne.MenuItem, 0)

		for _, u := range button.unlock.Users {
			user := u

			if user.UnpackStatus == unlocks.NotUnpacked {
				menuItem := fyne.NewMenuItem(
					"Unpack for "+user.ProfileName,
					func() {
						button.unlock.QueueUnpack(user)
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
		button.unlock.QueueUnpack(nil)
	}
}

func newUnpackButton(unlock *unlocks.Unlock) *unpackButton {
	button := &unpackButton{
		Button: widget.Button{
			Text: "Unpack",
			Icon: theme.FolderOpenIcon(),
		},

		unlock: unlock,
	}

	button.ExtendBaseWidget(button)
	return button
}
