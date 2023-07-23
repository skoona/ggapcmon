package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/skoona/ggapcmon/internal/commons"
	"image/color"
	"strconv"
)

func (v *viewProvider) OverviewPage() *fyne.Container {
	table := widget.NewTable(
		func() (int, int) { // length
			return len(v.prfHostKeys), 7
		},
		func() fyne.CanvasObject { // created
			i := widget.NewIcon(theme.StorageIcon())
			i.Hide()

			l := widget.NewLabel("0123456789")

			return container.NewHBox(i, l) // issue container minSize is 0
		},
		func(id widget.TableCellID, object fyne.CanvasObject) { // update
			// Row, Col
			host := v.cfg.HostByName(v.prfHostKeys[id.Row])
			switch id.Col {
			case 0: // State
				object.(*fyne.Container).Objects[0].(*widget.Icon).SetResource(commons.SknSelectResource(host.State))
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
				if host.Enabled {
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

	table.SetColumnWidth(0, 24)  // icon
	table.SetColumnWidth(1, 80)  // enabled
	table.SetColumnWidth(2, 104) // use tray
	table.SetColumnWidth(3, 132) // Name
	table.SetColumnWidth(4, 132) // Ip
	table.SetColumnWidth(5, 32)  // net period
	table.SetColumnWidth(6, 32)  // graph period

	rect := canvas.NewRectangle(color.Transparent)
	rect.StrokeWidth = 4
	rect.StrokeColor = theme.PrimaryColor()

	return container.NewPadded(rect, table)
}
