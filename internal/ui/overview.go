package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/skoona/ggapcmon/internal/commons"
	"image/color"
	"strings"
)

func (v *viewProvider) OverviewPage() *fyne.Container {
	table := widget.NewTable(
		func() (int, int) { // length, columns
			return len(v.prfHostKeys), 3
		},
		func() fyne.CanvasObject { // created
			i := widget.NewIcon(theme.StorageIcon())
			i.Hide()

			l := widget.NewRichTextFromMarkdown("")

			return container.NewHBox(i, l) // issue container minSize is 0
		},
		func(id widget.TableCellID, object fyne.CanvasObject) { // update
			// ICON - STATUS, ddd Outages, Last on dateString,
			//                LineV , DDD % percent Charge
			// Row, Col
			host := v.cfg.HostByName(v.prfHostKeys[id.Row])
			switch id.Col {
			case 0: // State
				object.(*fyne.Container).Objects[0].(*widget.Icon).SetResource(commons.SknSelectThemedResource(host.State))
				object.(*fyne.Container).Objects[0].(*widget.Icon).Resize(fyne.NewSize(40, 40))
				object.(*fyne.Container).Objects[0].Show()
				object.(*fyne.Container).Objects[1].Hide()

			case 1: // Enabled
				object.(*fyne.Container).Objects[0].Hide()
				object.(*fyne.Container).Objects[1].(*widget.RichText).ParseMarkdown("## " + strings.ToUpper(host.State))
				object.(*fyne.Container).Objects[1].Show()

			case 2: // descriptions
				z, _ := host.Bmaster.Get()
				if !host.Enabled {
					z = fmt.Sprintf("%s@%s\n\nno data available", host.Name, host.IpAddress)
				} else if z != "" {
					// network node
					a := z
					b, _ := host.Bcable.Get()
					z = fmt.Sprintf("%s@%s UPS host: %s\n\n Driver interface: %s",
						host.Name, host.IpAddress, a, b,
					)
				} else {
					x := z
					a, _ := host.Bloadpct.Get()
					b, _ := host.Bnumxfers.Get()
					c, _ := host.Bxonbatt.Get()
					d, _ := host.Blinev.Get()
					e, _ := host.Bbcharge.Get()
					z = fmt.Sprintf("## %s@%s Load %s\n\n%s Outages, Last on %s\n\n%s VAC, %s charge :%s",
						host.Name,
						host.IpAddress,
						a, b, c, d, e, x,
					)
				}

				object.(*fyne.Container).Objects[0].Hide()
				object.(*fyne.Container).Objects[1].(*widget.RichText).ParseMarkdown(z)
				object.(*fyne.Container).Objects[1].Refresh()
				object.(*fyne.Container).Objects[1].Show()
				object.(*fyne.Container).Refresh()

			default:
				object.(*fyne.Container).Objects[0].Hide()
				object.(*fyne.Container).Objects[1].(*widget.Label).SetText("Default")
				object.(*fyne.Container).Objects[1].Refresh()
				object.(*fyne.Container).Objects[1].Show()
			}
		},
	)

	table.SetColumnWidth(0, 56)  // icon
	table.SetColumnWidth(1, 96)  // status
	table.SetColumnWidth(2, 384) // description
	for idx := range v.cfg.Hosts() {
		table.SetRowHeight(idx, 72)
	}

	rect := canvas.NewRectangle(color.Transparent)
	rect.StrokeWidth = 4
	rect.StrokeColor = theme.PrimaryColor()

	v.overviewTable = table // allow external refresh

	return container.NewPadded(rect, table)
}
