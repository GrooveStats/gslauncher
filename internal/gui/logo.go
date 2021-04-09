package gui

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed groovestats-nobg.png
var logoData []byte

var groovestatsLogo = &fyne.StaticResource{
	StaticName:    "groovestats-nobg.png",
	StaticContent: logoData,
}
