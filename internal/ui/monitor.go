package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/skoona/ggapcmon/internal/commons"
	"github.com/skoona/ggapcmon/internal/entities"
	"github.com/skoona/ggapcmon/internal/interfaces"
	"github.com/skoona/sknlinechart"
	"image/color"
	"strconv"
	"strings"
	"time"
)

// monitor page

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
				object.(*fyne.Container).Objects[0].(*widget.Icon).SetResource(commons.SknSelectResource("unplugged"))
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
			//object.Refresh()
		},
	)

	table.SetColumnWidth(0, 24)  // icon
	table.SetColumnWidth(1, 96)  // enabled
	table.SetColumnWidth(2, 128) // use tray
	table.SetColumnWidth(3, 132) // Name
	table.SetColumnWidth(4, 132) // Ip
	table.SetColumnWidth(5, 32)  // net period
	table.SetColumnWidth(6, 32)  // graph period

	rect := canvas.NewRectangle(color.Transparent)
	rect.StrokeWidth = 4
	rect.StrokeColor = theme.PrimaryColor()

	return container.NewPadded(rect, table)
}

func (v *viewProvider) InfoPage(h entities.ApcHost) *fyne.Container {
	return container.NewGridWithColumns(4, widget.NewLabel(h.IpAddress))
}

func (v *viewProvider) MonitorPage() *fyne.Container {
	v.log.Println("HostKeys on Main Page: ", v.prfHostKeys, ", Hosts: ", v.prfHost)

	desc := canvas.NewText("Monitoring ", color.White)
	desc.Alignment = fyne.TextAlignCenter
	desc.TextStyle = fyne.TextStyle{Italic: true}
	desc.TextSize = 18

	hostTabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Overview", theme.ComputerIcon(), v.OverviewPage()),
		container.NewTabItemWithIcon("Glossary", theme.InfoIcon(), widget.NewLabel("Glossary")),
	)

	for _, host := range v.cfg.Hosts() {
		if host.Enabled {
			status := widget.NewMultiLineEntry()
			status.SetPlaceHolder("Status Page")
			status.TextStyle = fyne.TextStyle{Monospace: true}

			events := widget.NewMultiLineEntry()
			events.SetPlaceHolder("Events Page")
			events.TextStyle = fyne.TextStyle{Monospace: true}

			data := map[string][]*sknlinechart.ChartDatapoint{}
			chart, _ := sknlinechart.NewLineChart(host.Name, host.IpAddress, &data)

			info := v.InfoPage(host)

			tab := container.NewTabItemWithIcon(host.Name, theme.InfoIcon(),
				container.NewAppTabs(
					container.NewTabItemWithIcon("Chart", theme.ComputerIcon(), chart),
					container.NewTabItemWithIcon("Info", theme.InfoIcon(), info),
					container.NewTabItemWithIcon("Status", theme.WarningIcon(), status),
					container.NewTabItemWithIcon("Events", theme.StorageIcon(), events),
				),
			)
			hostTabs.Append(tab)

			go func(h entities.ApcHost, svc interfaces.Service, st *widget.Entry, ev *widget.Entry, chart sknlinechart.LineChart, info fyne.CanvasObject) {
				v.log.Println("ViewProvider::MonitorPage[", h.Name, "] BEGIN")
				rcvTuple := svc.MessageChannelByName(h.Name)
				var stBuild strings.Builder
				var evBuild strings.Builder
				var currentSt []string
				var currentEv []string
			pageExit:
				for {
					select {
					case <-v.ctx.Done():
						v.log.Println("ViewProvider::MonitorPage[", h.Name, "] fired:", v.ctx.Err().Error())
						break pageExit

					case msg := <-rcvTuple.Status:
						currentSt = msg
						stBuild.Reset()
						for idx, line := range currentSt {
							stBuild.WriteString(fmt.Sprintf("[%02d] %s\n", idx, line))
						}
						st.SetText(stBuild.String())

					case msg := <-rcvTuple.Events:
						currentEv = msg
						evBuild.Reset()
						for idx, line := range currentEv {
							evBuild.WriteString(fmt.Sprintf("[%02d] %s\n", idx, line))
						}
						ev.SetText(evBuild.String())
					default:
						var params map[string]string
						if len(currentSt) > 0 {
							info.(*fyne.Container).RemoveAll()
							params = svc.ParseStatus(currentSt)
							for k, v := range params {
								info.(*fyne.Container).Add(container.NewHBox(widget.NewLabel(k), widget.NewLabel(v)))
							}

							if len(currentSt) > 24 { // slaves have 21

								for k, v := range params {
									switch k {
									case "LINEV":
										d64, _ := strconv.ParseFloat(strings.TrimSpace(v), 32)
										point := sknlinechart.NewChartDatapoint(float32(d64), theme.ColorYellow, time.Now().Format(time.RFC1123))
										chart.ApplyDataPoint("LINEV", &point)
									case "LOADPCT":
										d64, _ := strconv.ParseFloat(strings.TrimSpace(v), 32)
										point := sknlinechart.NewChartDatapoint(float32(d64), theme.ColorBlue, time.Now().Format(time.RFC1123))
										chart.ApplyDataPoint("LOADPCT", &point)
									case "BCHARGE":
										d64, _ := strconv.ParseFloat(strings.TrimSpace(v), 32)
										point := sknlinechart.NewChartDatapoint(float32(d64), theme.ColorGreen, time.Now().Format(time.RFC1123))
										chart.ApplyDataPoint("BCHARGE", &point)
									case "TIMELEFT":
										d64, _ := strconv.ParseFloat(strings.TrimSpace(v), 32)
										point := sknlinechart.NewChartDatapoint(float32(d64), theme.ColorOrange, time.Now().Format(time.RFC1123))
										chart.ApplyDataPoint("TIMELEFT", &point)
									case "BATTV":
										d64, _ := strconv.ParseFloat(strings.TrimSpace(v), 32)
										point := sknlinechart.NewChartDatapoint(float32(d64), theme.ColorRed, time.Now().Format(time.RFC1123))
										chart.ApplyDataPoint("BATTV", &point)
									}
								}
							}

							currentSt = currentSt[:0]
						}
					}
				}
				v.log.Println("ViewProvider::MonitorPage[", h.Name, "] ENDED")
			}(host, v.service, status, events, chart, info)

		}
	}

	page := container.NewBorder(
		desc,
		nil,
		nil,
		nil,
		hostTabs)
	return page
}
