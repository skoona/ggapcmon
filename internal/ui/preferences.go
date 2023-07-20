package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/skoona/ggapcmon/internal/commons"
	"image/color"
	"strconv"
)

// settings page
func (v *viewProvider) PrefsPage() *fyne.Container {
	v.prefsHostKeys = v.cfg.HostKeys()
	v.prefsHost = v.cfg.HostByName(v.prefsHostKeys[0])

	sDesc := canvas.NewText("Selected Host", color.White)
	sDesc.Alignment = fyne.TextAlignLeading
	sDesc.TextStyle = fyne.TextStyle{Italic: true}
	sDesc.TextSize = 24
	desc := container.NewPadded(sDesc)

	tdesc := canvas.NewText("Hosts", color.White)
	tdesc.Alignment = fyne.TextAlignLeading
	tdesc.TextStyle = fyne.TextStyle{Italic: true}
	tdesc.TextSize = 24
	tDesc := container.NewPadded(tdesc)

	dHost := widget.NewEntry()
	dHost.SetPlaceHolder("10.100.1.3:3551")
	dName := widget.NewEntry()
	dName.SetPlaceHolder("VServ")
	nPeriod := widget.NewEntry()
	nPeriod.SetPlaceHolder("15")
	gPeriod := widget.NewEntry()
	gPeriod.SetPlaceHolder("30")

	enable := widget.NewCheck("", func(onOff bool) {
		if onOff {
			h := v.cfg.HostByName(v.prefsHost.Name)
			h.Enabled = true
			v.cfg.Apply(h)
		} else {
			h := v.cfg.HostByName(v.prefsHost.Name)
			h.Enabled = false
			v.cfg.Apply(h)
		}
	})
	//enable.SetChecked(commons.IsInfluxDBEnabled())

	trayIcon := widget.NewCheck("", func(onOff bool) {
		if onOff {
			h := v.cfg.HostByName(v.prefsHost.Name)
			h.TrayIcon = true
			v.cfg.Apply(h)
		} else {
			h := v.cfg.HostByName(v.prefsHost.Name)
			h.TrayIcon = false
			v.cfg.Apply(h)
		}
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
		v.cfg.Apply(v.prefsHost)
		fmt.Println("Form submitted: restart for effect", v.prefsHost)
	}

	table := widget.NewTable(
		func() (int, int) { // length
			return len(v.prefsHostKeys), 7
		},
		func() fyne.CanvasObject { // created
			return container.NewPadded()
		},
		func(id widget.TableCellID, object fyne.CanvasObject) { // update
			var elem fyne.CanvasObject
			// Row, Col
			host := v.cfg.HostByName(v.prefsHostKeys[id.Row])
			switch id.Col {
			case 0: // State
				elem = commons.SknSelectImage("unplugged")
			case 1: // Enabled
				elem = widget.NewCheck("enabled", func(b bool) {
					v.log.Println("Enable: Row:Col", id.Row, ":", id.Col)
				})
				elem.(*widget.Check).SetChecked(host.Enabled)
			case 2: // Tray
				elem = widget.NewCheck("use trayIcon", func(b bool) {
					v.log.Println("Tray: Row:Col", id.Row, ":", id.Col)
				})
				elem.(*widget.Check).SetChecked(host.TrayIcon)
			case 3: // Name
				elem = widget.NewLabel(host.Name)
			case 4: // IP
				elem = widget.NewLabel(host.IpAddress)
			case 5: // Network
				elem = widget.NewLabel(strconv.Itoa(int(host.NetworkSamplePeriod)))
			case 6: // Graph
				elem = widget.NewLabel(strconv.Itoa(int(host.GraphingSamplePeriod)))
			default:
				elem = widget.NewLabel("Default")
				elem.(*widget.Label).TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
			}
			object.(*fyne.Container).Add(elem)
		},
	)
	table.OnSelected = func(id widget.TableCellID) {
		v.log.Println("Selected: ", id.Row, ":", id.Col, ", Host: ", v.cfg.HostByName(v.prefsHostKeys[id.Row]))
	}
	//table.SetRowHeight(0, 24)
	//table.SetRowHeight(1, 24)
	//table.SetColumnWidth(0, 32)
	//table.SetColumnWidth(1, 32)
	//table.SetColumnWidth(2, 128)
	//table.SetColumnWidth(3, 128)
	//table.SetColumnWidth(4, 128)
	//table.SetColumnWidth(5, 32)
	//table.SetColumnWidth(6, 32)

	page := container.NewGridWithColumns(1,
		settings.NewSettings().LoadAppearanceScreen(v.mainWindow),
		container.NewBorder(
			desc,
			nil,
			nil,
			nil,
			form,
		),

		container.NewBorder(
			tDesc,
			container.NewHBox(
				widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
					v.log.Println("Refresh clicked")
				}),
				widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
					v.log.Println("Add clicked")
				}),
				widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), func() {
					v.log.Println("Remove clicked")
				}),
			),
			nil,
			nil,
			container.NewMax(table),
		),
	)
	return page
}
