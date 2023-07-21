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
	"time"
)

// settings page
func (v *viewProvider) PrefsPage() *fyne.Container {
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
	dHost.SetText(v.prfHost.IpAddress)

	dName := widget.NewEntry()
	dName.SetText(v.prfHost.Name)

	z := strconv.Itoa(int(v.prfHost.NetworkSamplePeriod))
	nPeriod := widget.NewEntry()
	nPeriod.SetText(z)

	z = strconv.Itoa(int(v.prfHost.GraphingSamplePeriod))
	gPeriod := widget.NewEntry()
	gPeriod.SetText(z)

	enable := widget.NewCheck("", func(onOff bool) {
		if onOff {
			h := v.cfg.HostByName(v.prfHost.Name)
			h.Enabled = true
			v.cfg.Apply(h)
		} else {
			h := v.cfg.HostByName(v.prfHost.Name)
			h.Enabled = false
			v.cfg.Apply(h)
		}
	})
	enable.SetChecked(v.prfHost.Enabled)

	trayIcon := widget.NewCheck("", func(onOff bool) {
		if onOff {
			h := v.cfg.HostByName(v.prfHost.Name)
			h.TrayIcon = true
			v.cfg.Apply(h)
		} else {
			h := v.cfg.HostByName(v.prfHost.Name)
			h.TrayIcon = false
			v.cfg.Apply(h)
		}
	})
	trayIcon.SetChecked(v.prfHost.TrayIcon)

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
		//if reflect.TypeOf(v.prfHost).String() == "entities.ApcHost" {
		v.prfHost.Name = dName.Text
		v.prfHost.IpAddress = dHost.Text
		x, _ := strconv.Atoi(nPeriod.Text)
		v.prfHost.NetworkSamplePeriod = time.Duration(x)
		x, _ = strconv.Atoi(gPeriod.Text)
		v.prfHost.GraphingSamplePeriod = time.Duration(x)
		v.prfHost.TrayIcon = trayIcon.Checked
		v.prfHost.Enabled = enable.Checked

		v.cfg.Apply(v.prfHost).Save()
		v.ShowPrefsPage()
		fmt.Println("Form submitted: restart for effect", v.prfHost)
		//}
	}

	table := widget.NewTable(
		func() (int, int) { // length
			return len(v.prfHostKeys), 7
		},
		func() fyne.CanvasObject { // created
			i := widget.NewIcon(theme.StorageIcon())
			i.Hide()

			c := widget.NewCheck("0123456789", nil)
			c.Hide()

			l := widget.NewLabel("0123456789")

			return container.NewHBox(i, c, l) // issue container minSize is 0
		},
		func(id widget.TableCellID, object fyne.CanvasObject) { // update
			// Row, Col
			host := v.cfg.HostByName(v.prfHostKeys[id.Row])
			switch id.Col {
			case 0: // State
				object.(*fyne.Container).Objects[0].(*widget.Icon).SetResource(commons.SknSelectResource("unplugged"))
				object.(*fyne.Container).Objects[1].Hide()
				object.(*fyne.Container).Objects[2].Hide()
				object.(*fyne.Container).Objects[0].Refresh()
				object.(*fyne.Container).Objects[0].Show()

			case 1: // Enabled
				object.(*fyne.Container).Objects[0].Hide()
				object.(*fyne.Container).Objects[1].(*widget.Check).Text = "enabled"
				object.(*fyne.Container).Objects[1].(*widget.Check).SetChecked(host.Enabled)
				//object.(*fyne.Container).Objects[1].(*widget.Check).OnChanged = func(b bool) {
				//	v.log.Println("Enable: Row:Col", id.Row, ":", id.Col)
				//}
				object.(*fyne.Container).Objects[2].Hide()
				object.(*fyne.Container).Objects[1].Refresh()
				object.(*fyne.Container).Objects[1].Show()

			case 2: // Tray
				object.(*fyne.Container).Objects[0].Hide()
				object.(*fyne.Container).Objects[1].(*widget.Check).Text = "use trayIcon"
				object.(*fyne.Container).Objects[1].(*widget.Check).SetChecked(host.TrayIcon)
				//object.(*fyne.Container).Objects[1].(*widget.Check).OnChanged = func(b bool) {
				//	v.log.Println("TryIcon: Row:Col", id.Row, ":", id.Col)
				//}
				object.(*fyne.Container).Objects[2].Hide()
				object.(*fyne.Container).Objects[1].Refresh()
				object.(*fyne.Container).Objects[1].Show()

			case 3: // Name
				object.(*fyne.Container).Objects[0].Hide()
				object.(*fyne.Container).Objects[1].Hide()
				object.(*fyne.Container).Objects[2].(*widget.Label).SetText(host.Name)
				object.(*fyne.Container).Objects[2].Refresh()
				object.(*fyne.Container).Objects[2].Show()

			case 4: // IP
				object.(*fyne.Container).Objects[0].Hide()
				object.(*fyne.Container).Objects[1].Hide()
				object.(*fyne.Container).Objects[2].(*widget.Label).SetText(host.IpAddress)
				object.(*fyne.Container).Objects[2].Refresh()
				object.(*fyne.Container).Objects[2].Show()

			case 5: // Network
				object.(*fyne.Container).Objects[0].Hide()
				object.(*fyne.Container).Objects[1].Hide()
				object.(*fyne.Container).Objects[2].(*widget.Label).SetText(strconv.Itoa(int(host.NetworkSamplePeriod)))
				object.(*fyne.Container).Objects[2].Refresh()
				object.(*fyne.Container).Objects[2].Show()

			case 6: // Graph
				object.(*fyne.Container).Objects[0].Hide()
				object.(*fyne.Container).Objects[1].Hide()
				object.(*fyne.Container).Objects[2].(*widget.Label).SetText(strconv.Itoa(int(host.GraphingSamplePeriod)))
				object.(*fyne.Container).Objects[2].Refresh()
				object.(*fyne.Container).Objects[2].Show()

			default:
				object.(*fyne.Container).Objects[0].Hide()
				object.(*fyne.Container).Objects[1].Hide()
				object.(*fyne.Container).Objects[2].(*widget.Label).SetText("Default")
				object.(*fyne.Container).Objects[2].Refresh()
				object.(*fyne.Container).Objects[2].Show()
			}
			//object.Refresh()
		},
	)
	table.OnSelected = func(id widget.TableCellID) {
		v.prfHost = v.cfg.HostByName(v.prfHostKeys[id.Row])

		dName.Text = v.prfHost.Name
		dHost.Text = v.prfHost.IpAddress
		z := strconv.Itoa(int(v.prfHost.NetworkSamplePeriod))
		nPeriod.Text = z
		z = strconv.Itoa(int(v.prfHost.GraphingSamplePeriod))
		gPeriod.Text = z
		trayIcon.Checked = v.prfHost.TrayIcon
		enable.Checked = v.prfHost.Enabled

		form.Refresh()
		v.log.Println("Selected: ", id.Row, ":", id.Col, ", Host: ", v.cfg.HostByName(v.prfHostKeys[id.Row]))
	}
	table.SetColumnWidth(0, 24)  // icon
	table.SetColumnWidth(1, 96)  // enabled
	table.SetColumnWidth(2, 128) // use tray
	table.SetColumnWidth(3, 132) // Name
	table.SetColumnWidth(4, 132) // Ip
	table.SetColumnWidth(5, 32)  // net period
	table.SetColumnWidth(6, 32)  // graph period

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
					v.ShowPrefsPage()
					v.log.Println("Refresh clicked")
				}),
				widget.NewButtonWithIcon("add selected", theme.ContentAddIcon(), func() {
					form.OnSubmit()
					v.prefsAddAction()
					v.log.Println("Add clicked")
				}),
				widget.NewButtonWithIcon("del selected", theme.ContentRemoveIcon(), func() {
					// v.prfHost will be deleted
					v.prefsDelAction()
					v.log.Println("Remove clicked")
				}),
			),
			nil,
			nil,
			table,
		),
	)
	return page
}
