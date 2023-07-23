package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

func (v *viewProvider) Performance(status map[string]string) *fyne.Container {
	desc := canvas.NewText("Performance Summary", theme.PrimaryColor())
	desc.Alignment = fyne.TextAlignCenter
	desc.TextStyle = fyne.TextStyle{Italic: true}
	desc.TextSize = 18

	frame := canvas.NewRectangle(color.Transparent)
	frame.StrokeWidth = 6
	frame.StrokeColor = theme.PlaceHolderColor()

	items := container.New(layout.NewFormLayout())

	titleBorder := container.NewPadded(
		frame,
		container.NewBorder(
			container.NewPadded(
				canvas.NewRectangle(theme.PlaceHolderColor()),
				desc,
			),
			nil,
			nil,
			nil,
			items,
		),
	)

	st := status["SELFTEST"]
	lbl := widget.NewLabel("Selftest running")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabel(st))

	st = status["NUMXFERS"]
	lbl = widget.NewLabel("Number of transfers")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabel(st))

	st = status["LASTXFER"]
	lbl = widget.NewLabel("Reason last transfer")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabel(st))

	st = status["XONBATT"]
	lbl = widget.NewLabel("Last transfer to battery")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabel(st))

	st = status["XOFFBATT"]
	lbl = widget.NewLabel("Last transfer off battery")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabel(st))

	st = status["TONBATT"]
	lbl = widget.NewLabel("Time on battery")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabel(st))

	st = status["CUMONBATT"]
	lbl = widget.NewLabel("Cummulative on battery")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabel(st))

	return titleBorder
}
func (v *viewProvider) Metrics(status map[string]string) *fyne.Container {
	desc := canvas.NewText("UP Metrics", theme.PrimaryColor())
	desc.Alignment = fyne.TextAlignCenter
	desc.TextStyle = fyne.TextStyle{Italic: true}
	desc.TextSize = 18

	frame := canvas.NewRectangle(color.Transparent)
	frame.StrokeWidth = 6
	frame.StrokeColor = theme.PlaceHolderColor()

	items := container.New(layout.NewFormLayout())

	titleBorder := container.NewPadded(
		frame,
		container.NewBorder(
			container.NewPadded(
				canvas.NewRectangle(theme.PlaceHolderColor()),
				desc,
			),
			nil,
			nil,
			nil,
			items,
		),
	)

	st := status["LINEV"]
	lbl := widget.NewLabel("Utility line")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabel(st))

	st = status["BATTV"]
	lbl = widget.NewLabel("Battery DC")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabel(st))

	st = status["BCHARGE"]
	lbl = widget.NewLabel("Percent battery charge")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabel(st))

	st = status["LOADPCT"]
	lbl = widget.NewLabel("Percent load capacity")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabel(st))

	st = status["TIMELEFT"]
	lbl = widget.NewLabel("Minutes remaining")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabel(st))

	return titleBorder
}
func (v *viewProvider) Software(status map[string]string) *fyne.Container {
	desc := canvas.NewText("Software Information ", theme.PrimaryColor())
	desc.Alignment = fyne.TextAlignCenter
	desc.TextStyle = fyne.TextStyle{Italic: true}
	desc.TextSize = 18

	frame := canvas.NewRectangle(color.Transparent)
	frame.StrokeWidth = 6
	frame.StrokeColor = theme.PlaceHolderColor()

	items := container.New(layout.NewFormLayout())

	titleBorder := container.NewPadded(
		frame,
		container.NewBorder(
			container.NewPadded(
				canvas.NewRectangle(theme.PlaceHolderColor()),
				desc,
			),
			nil,
			nil,
			nil,
			items,
		),
	)

	st := status["VERSION"]
	lbl := widget.NewLabel("APCUPSD version")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabel(st))

	st = status["UPSNAME"]
	lbl = widget.NewLabel("Monitored UPS name")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabel(st))

	st = status["CABLE"]
	lbl = widget.NewLabel("Cable driver type")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabel(st))

	st = status["UPSMODE"]
	lbl = widget.NewLabel("Configuration mode")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabel(st))

	st = status["STARTTIME"]
	lbl = widget.NewLabel("Last started")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabel(st))

	st = status["STATUS"]
	lbl = widget.NewLabel("UPS state")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabel(st))

	return titleBorder
}
func (v *viewProvider) Product(status map[string]string) *fyne.Container {
	desc := canvas.NewText("Product Information ", theme.PrimaryColor())
	desc.Alignment = fyne.TextAlignCenter
	desc.TextStyle = fyne.TextStyle{Italic: true}
	desc.TextSize = 18

	frame := canvas.NewRectangle(color.Transparent)
	frame.StrokeWidth = 6
	frame.StrokeColor = theme.PlaceHolderColor()

	items := container.New(layout.NewFormLayout())

	titleBorder := container.NewPadded(
		frame,
		container.NewBorder(
			container.NewPadded(
				canvas.NewRectangle(theme.PlaceHolderColor()),
				desc,
			),
			nil,
			nil,
			nil,
			items,
		),
	)

	st := status["MODEL"]
	lbl := widget.NewLabel("Device model")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabel(st))

	st = status["SERIALNO"]
	lbl = widget.NewLabel("Serial number")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabel(st))

	st = status["MANDATE"]
	lbl = widget.NewLabel("Manufacture date")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabel(st))

	st = status["FIRMWARE"]
	lbl = widget.NewLabel("Firmware")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabel(st))

	st = status["BATTDATE"]
	lbl = widget.NewLabel("Battery date")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabel(st))

	st = status["ITEMP"]
	lbl = widget.NewLabel("Internal temp")
	lbl.Alignment = fyne.TextAlignTrailing
	items.Add(lbl)
	items.Add(widget.NewLabel(st))

	return titleBorder
}

func (v *viewProvider) DetailPage(params chan map[string]string) *fyne.Container {
	page := container.NewGridWithColumns(2)

	go func() {
		for status := range params {
			page.RemoveAll()
			page.Add(v.Performance(status))
			page.Add(v.Metrics(status))
			page.Add(v.Software(status))
			page.Add(v.Product(status))
		}
	}()

	return page
}
