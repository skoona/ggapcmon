package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

// settings page
func (v *viewProvider) PrefsPage() *fyne.Container {
	desc := canvas.NewText("Configuration", color.White)
	desc.Alignment = fyne.TextAlignLeading
	desc.TextStyle = fyne.TextStyle{Italic: true}
	desc.TextSize = 24

	dHost := widget.NewEntry()
	dHost.SetPlaceHolder("10.100.1.3:3551")
	dName := widget.NewEntry()
	dName.SetPlaceHolder("VServ")
	nPeriod := widget.NewEntry()
	nPeriod.SetPlaceHolder("15")
	gPeriod := widget.NewEntry()
	gPeriod.SetPlaceHolder("30")

	enable := widget.NewCheck("", func(onOff bool) {
		//if onOff {
		//	commons.SetEnableInfluxDB("true")
		//} else {
		//	commons.SetEnableInfluxDB("false")
		//}
	})
	//enable.SetChecked(commons.IsInfluxDBEnabled())

	trayIcon := widget.NewCheck("", func(onOff bool) {
		//if onOff {
		//	commons.SetEnableDebugMode("true")
		//} else {
		//	commons.SetEnableDebugMode("false")
		//}
	})
	//trayIcon.SetChecked(commons.IsDebugMode())

	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "host name", Widget: dName},
			{Text: "host URI", Widget: dHost},
			{Text: "graph averaging count", Widget: gPeriod},
			{Text: "network sampling period", Widget: nPeriod},
			{Text: "use tray icon", Widget: trayIcon},
			{Text: "is enabled", Widget: enable},
		},
		SubmitText: "Apply",
	}
	form.OnSubmit = func() { // optional, handle form submission

		fmt.Println("Form submitted: restart for effect")
	}

	table := widget.NewTable(
		func() (int, int) {
			return len(v.cfg.Hosts()), 6
		},
		func() fyne.CanvasObject {
			return container.NewCenter()
		},
		func(id widget.TableCellID, object fyne.CanvasObject) {

		},
	)

	page := container.NewVBox(
		settings.NewSettings().LoadAppearanceScreen(v.mainWindow),
		container.NewBorder(desc, nil, nil, nil, form, table),
	)
	return page
}
