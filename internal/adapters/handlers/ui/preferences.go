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
	"github.com/skoona/ggapcmon/internal/core/domain"
	"image/color"
	"strconv"
	"time"
)

// PreferencesPage manages application settings
func (v *viewProvider) PreferencesPage() *fyne.Container {
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

	table := widget.NewTable(
		func() (int, int) { // length
			return len(v.prfHostKeys) + 1, 7
		},
		func() fyne.CanvasObject { // created
			//i := widget.NewIcon(theme.StorageIcon())
			i := commons.SknSelectThemedImage("charging")
			i.Hide()
			l := widget.NewLabel("0123456789")
			return container.NewHBox(i, l) // issue container minSize is 0
		},
		func(id widget.TableCellID, object fyne.CanvasObject) { // update
			// ICON - STATUS, ddd Outages, Last on dateString,
			//                LineV , DDD % percent Charge
			// Row, Col
			if id.Row == 0 { // headers
				object.(*fyne.Container).Objects[0].Hide()
				switch id.Col {
				case 0:
					object.(*fyne.Container).Objects[1].(*widget.Label).SetText("State")
					object.(*fyne.Container).Objects[1].Show()
				case 1:
					object.(*fyne.Container).Objects[1].(*widget.Label).SetText("Enabled")
					object.(*fyne.Container).Objects[1].Show()
				case 2:
					object.(*fyne.Container).Objects[1].(*widget.Label).SetText("Tray")
					object.(*fyne.Container).Objects[1].Show()
				case 3:
					object.(*fyne.Container).Objects[1].(*widget.Label).SetText("Name")
					object.(*fyne.Container).Objects[1].Show()
				case 4:
					object.(*fyne.Container).Objects[1].(*widget.Label).SetText("IpAddress")
					object.(*fyne.Container).Objects[1].Show()
				case 5:
					object.(*fyne.Container).Objects[1].(*widget.Label).SetText("NetPeriod")
					object.(*fyne.Container).Objects[1].Show()
				case 6:
					object.(*fyne.Container).Objects[1].(*widget.Label).SetText("GraphPeriod")
					object.(*fyne.Container).Objects[1].Show()
				}
				return
			}
			// Row, Col
			host := v.cfg.HostById(v.prfHostKeys[id.Row-1])
			switch id.Col {
			case 0: // State
				object.(*fyne.Container).Objects[0].(*canvas.Image).Resource = commons.SknSelectThemedResource(host.State)
				object.(*fyne.Container).Objects[0].(*canvas.Image).Resize(fyne.NewSize(40, 40))
				object.(*fyne.Container).Objects[1].Hide()
				object.(*fyne.Container).Objects[0].Refresh()
				object.(*fyne.Container).Objects[0].Show()

			case 1: // Enabled
				label := "disabled"
				if host.Enabled {
					label = "enabled"
				}
				object.(*fyne.Container).Objects[0].Hide()
				object.(*fyne.Container).Objects[1].(*widget.Label).SetText(label)
				object.(*fyne.Container).Objects[1].Refresh()
				object.(*fyne.Container).Objects[1].Show()

			case 2: // Tray
				label := "no trayIcon"
				if host.TrayIcon {
					label = "use trayIcon"
				}
				object.(*fyne.Container).Objects[0].Hide()
				object.(*fyne.Container).Objects[1].(*widget.Label).SetText(label)
				object.(*fyne.Container).Objects[1].Refresh()
				object.(*fyne.Container).Objects[1].Show()

			case 3: // Name
				object.(*fyne.Container).Objects[0].Hide()
				object.(*fyne.Container).Objects[1].(*widget.Label).SetText(host.Name)
				object.(*fyne.Container).Objects[1].Refresh()
				object.(*fyne.Container).Objects[1].Show()

			case 4: // IP
				object.(*fyne.Container).Objects[0].Hide()
				object.(*fyne.Container).Objects[1].(*widget.Label).SetText(host.IpAddress)
				object.(*fyne.Container).Objects[1].Refresh()
				object.(*fyne.Container).Objects[1].Show()

			case 5: // Network
				object.(*fyne.Container).Objects[0].Hide()
				object.(*fyne.Container).Objects[1].(*widget.Label).SetText(strconv.Itoa(int(host.NetworkSamplePeriod)))
				object.(*fyne.Container).Objects[1].Refresh()
				object.(*fyne.Container).Objects[1].Show()

			case 6: // Graph
				object.(*fyne.Container).Objects[0].Hide()
				object.(*fyne.Container).Objects[1].(*widget.Label).SetText(strconv.Itoa(int(host.GraphingSamplePeriod)))
				object.(*fyne.Container).Objects[1].Refresh()
				object.(*fyne.Container).Objects[1].Show()

			default:
				object.(*fyne.Container).Objects[0].Hide()
				object.(*fyne.Container).Objects[1].(*widget.Label).SetText("Default")
				object.(*fyne.Container).Objects[1].Refresh()
				object.(*fyne.Container).Objects[1].Show()
			}
		},
	)

	dHost := widget.NewEntry()
	dHost.SetText(v.prfHost.IpAddress)

	dId := widget.NewEntry()
	dId.SetText(v.prfHost.Id)
	dId.Disable()

	dName := widget.NewEntry()
	dName.SetText(v.prfHost.Name)

	z := strconv.Itoa(int(v.prfHost.NetworkSamplePeriod))
	nPeriod := widget.NewEntry()
	nPeriod.SetText(z)

	z = strconv.Itoa(int(v.prfHost.GraphingSamplePeriod))
	gPeriod := widget.NewEntry()
	gPeriod.SetText(z)

	enable := widget.NewCheck("", nil)
	enable.SetChecked(v.prfHost.Enabled)

	trayIcon := widget.NewCheck("", nil)
	trayIcon.SetChecked(v.prfHost.TrayIcon)

	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "host name", Widget: dName},
			{Text: "host URI", Widget: dHost},
			{Text: "graph averaging count", Widget: gPeriod},
			{Text: "network sampling period", Widget: nPeriod},
			{Text: "use tray icon", Widget: trayIcon},
			{Text: "is enabled", Widget: enable},
			{Text: "ID", Widget: dId},
		},
		SubmitText: "Apply",
	}
	form.OnSubmit = func() { // Apply optional, handle form submission
		if dId.Text == "" {
			nx, _ := strconv.Atoi(nPeriod.Text)
			gx, _ := strconv.Atoi(gPeriod.Text)
			v.prfHost = domain.NewApcHost(
				dName.Text,
				dHost.Text,
				time.Duration(nx),
				time.Duration(gx),
				enable.Checked,
				trayIcon.Checked,
			)
		} else {
			v.prfHost.Name = dName.Text
			v.prfHost.IpAddress = dHost.Text
			x, _ := strconv.Atoi(nPeriod.Text)
			v.prfHost.NetworkSamplePeriod = time.Duration(x)
			x, _ = strconv.Atoi(gPeriod.Text)
			v.prfHost.GraphingSamplePeriod = time.Duration(x)
			v.prfHost.TrayIcon = trayIcon.Checked
			v.prfHost.Enabled = enable.Checked
		}
		v.cfg.Apply(v.prfHost).Save()
		table.Refresh()
		v.prfStatusLine.SetText(fmt.Sprintf("Form submitted for host:%s, restart for effect", v.prfHost.Name))
	}

	table.OnSelected = func(id widget.TableCellID) {
		if id.Row == 0 {
			v.prfStatusLine.SetText(fmt.Sprintf("header column: %d selected", id.Col))
			return
		}

		if id.Row-1 > len(v.prfHostKeys) {
			v.prfHostKeys = v.cfg.HostKeys()
		}
		v.prfHost = v.cfg.HostById(v.prfHostKeys[id.Row-1])

		dId.SetText(v.prfHost.Id)
		dName.Text = v.prfHost.Name
		dHost.Text = v.prfHost.IpAddress
		z := strconv.Itoa(int(v.prfHost.NetworkSamplePeriod))
		nPeriod.Text = z
		z = strconv.Itoa(int(v.prfHost.GraphingSamplePeriod))
		gPeriod.Text = z
		trayIcon.Checked = v.prfHost.TrayIcon
		enable.Checked = v.prfHost.Enabled

		form.Refresh()
		v.prfStatusLine.SetText(fmt.Sprintf("Selected row:%d, col:%d, for host:%s", id.Row-1, id.Col, v.cfg.HostById(v.prfHostKeys[id.Row-1]).Name))
	}
	table.SetColumnWidth(0, 40)  // icon
	table.SetColumnWidth(1, 80)  // enabled
	table.SetColumnWidth(2, 104) // use tray
	table.SetColumnWidth(3, 132) // Name
	table.SetColumnWidth(4, 132) // Ip
	table.SetColumnWidth(5, 80)  // net period
	table.SetColumnWidth(6, 80)  // graph period

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
				container.NewHBox(
					widget.NewButtonWithIcon("add", theme.ContentAddIcon(), func() {
						dId.SetText("")
						form.OnSubmit()
						v.prfHostKeys = v.cfg.HostKeys()
						v.prfHost = v.cfg.HostById(v.prfHost.Id)
						table.Refresh()
						v.prfStatusLine.SetText("Host " + v.prfHost.Name + " was added")
					}),
					widget.NewButtonWithIcon("del", theme.ContentRemoveIcon(), func() {
						h := v.prfHost
						v.cfg.Remove(h.Id)
						v.prfHostKeys = v.cfg.HostKeys()
						v.prfHost = v.cfg.HostById(v.prfHostKeys[0])
						table.Refresh()
						v.prfStatusLine.SetText("Host " + h.Name + " was removed")
					}),
					widget.NewButtonWithIcon("test", theme.QuestionIcon(), func() {
						_ = v.verifyHostConnection()
						v.prfHostKeys = v.cfg.HostKeys()
						v.prfHost = v.cfg.HostById(v.prfHost.Id)
						table.Refresh()
					}),
				),
				v.prfStatusLine,
			),
			nil,
			nil,
			table,
		),
	)
	return page
}
