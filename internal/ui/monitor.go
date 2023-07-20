package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image/color"
)

// monitor page

func (v *viewProvider) MonitorPage() *fyne.Container {
	desc := canvas.NewText("Monitor", color.White)
	desc.Alignment = fyne.TextAlignCenter
	desc.TextStyle = fyne.TextStyle{Italic: true}
	desc.TextSize = 24

	place := canvas.NewText("StandBy: Page Under Construction", color.White)
	place.Alignment = fyne.TextAlignCenter
	place.TextStyle = fyne.TextStyle{Italic: true}
	place.TextSize = 24

	page := container.NewBorder(desc, nil, nil, nil, place)
	return page
}
