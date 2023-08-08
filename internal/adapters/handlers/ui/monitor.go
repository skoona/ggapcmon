package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/skoona/ggapcmon/internal/commons"
	"github.com/skoona/ggapcmon/internal/core/domain"
	"github.com/skoona/ggapcmon/internal/core/ports"
	"github.com/skoona/sknlinechart"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// MonitorPage primary application page
func (v *viewProvider) MonitorPage() *fyne.Container {
	hostTabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Known UPSs", theme.ComputerIcon(), v.OverviewPage()),
		container.NewTabItemWithIcon("Glossary", theme.InfoIcon(), v.GlossaryPage()),
	)

	for _, host := range v.cfg.Hosts() {
		if host.Enabled {
			v.bondedUpsStatus[host.Name] = domain.NewUpsStatusValueBindings(host)
			status := widget.NewMultiLineEntry()
			status.SetPlaceHolder("StandBy: no status has been received.")
			status.TextStyle = fyne.TextStyle{Monospace: true}

			events := widget.NewMultiLineEntry()
			events.SetPlaceHolder("StandBy: no events have been received.")
			events.TextStyle = fyne.TextStyle{Monospace: true}

			// GraphSamplingPeriod for Charts
			v.chartPageData[host.Name] = map[string]ports.GraphPointSmoothing{}
			for _, key := range v.chartKeys {
				intf := domain.NewGraphAverage(host.Name, key, host.GraphingSamplePeriod)
				v.chartPageData[host.Name][key] = intf
			}
			// for chart page updates
			data := map[string][]*sknlinechart.ChartDatapoint{}
			chart, err := sknlinechart.New(host.Name, "", 10, &data)
			if err != nil {
				dialog.ShowError(err, v.mainWindow)
				commons.ShutdownSignals <- syscall.SIGINT
			}
			chart.SetBottomLeftLabel(host.Name + "@" + host.IpAddress + " is " + host.State)

			// for details page updates
			knowledge := make(chan map[string]string, 16)

			tab := container.NewTabItemWithIcon(host.Name, theme.InfoIcon(),
				container.NewAppTabs(
					container.NewTabItemWithIcon("History", theme.HistoryIcon(), chart),
					container.NewTabItemWithIcon("Detailed", theme.VisibilityIcon(), container.NewScroll(v.DetailPage(knowledge, v.bondedUpsStatus[host.Name]))),
					container.NewTabItemWithIcon("Status", theme.ListIcon(), status),
					container.NewTabItemWithIcon("Events", theme.WarningIcon(), events),
				),
			)
			hostTabs.Append(tab)

			v.handleUpdatesForMonitorPage(host, v.service, status, events, chart, knowledge)
		}
	}

	return container.NewPadded(hostTabs)
}

// handleUpdatesForMonitorPage does exactly that
func (v *viewProvider) handleUpdatesForMonitorPage(host *domain.ApcHost, service ports.Service, status *widget.Entry, events *widget.Entry, chart sknlinechart.LineChart, kChan chan map[string]string) {
	go func(h *domain.ApcHost, svc ports.Service, st *widget.Entry, ev *widget.Entry, chart sknlinechart.LineChart, knowledge chan map[string]string) {
		commons.DebugLog("ViewProvider::HandleUpdatesForMonitorPage[", h.Name, "] BEGIN")
		rcvTuple := svc.MessageChannelByName(h.Name)
		var stBuild strings.Builder
		var evBuild strings.Builder
		var currentSt []string
		var currentEv []string
	pageExit:
		for {
			select {
			case <-v.ctx.Done():
				close(knowledge) // detail pages
				v.bondedUpsStatus[h.Name].UnbindUpsData()
				commons.DebugLog("ViewProvider::HandleUpdatesForMonitorPage[", h.Name, "] fired:", v.ctx.Err().Error())
				break pageExit

			case msg := <-rcvTuple.Status:
				currentSt = msg
				stBuild.Reset()
				for idx, line := range currentSt {
					stBuild.WriteString(fmt.Sprintf("[%02d] %s\n", idx, line))
				}
				st.SetText(stBuild.String())
				st.Refresh()

			case msg := <-rcvTuple.Events:
				currentEv = msg
				evBuild.Reset()
				for idx, line := range currentEv {
					evBuild.WriteString(fmt.Sprintf("[%02d] %s\n", idx, line))
				}
				ev.SetText(evBuild.String())
				ev.Refresh()

			default:
				var params map[string]string

				if len(currentSt) > 0 {
					params = svc.ParseStatus(currentSt)

					for k, vv := range params {
						floatStr := strings.Split(vv, " ")
						floatStr[0] = strings.TrimSpace(floatStr[0])
						// gapcmon charted: LINEV, LOADPCT, BCHARGE, CUMONBATT, TIMELEFT
						switch k {
						case "LINEV":
							d64, _ := strconv.ParseFloat(strings.TrimSpace(floatStr[0]), 32)
							d64 = v.chartPageData[h.Name][k].AddValue(d64)
							point := sknlinechart.NewChartDatapoint(float32(d64), theme.ColorYellow, time.Now().Format(time.RFC1123))
							chart.ApplyDataPoint("LINEV", &point)
						case "LOADPCT":
							d64, _ := strconv.ParseFloat(strings.TrimSpace(floatStr[0]), 32)
							d64 = v.chartPageData[h.Name][k].AddValue(d64)
							point := sknlinechart.NewChartDatapoint(float32(d64), theme.ColorBlue, time.Now().Format(time.RFC1123))
							chart.ApplyDataPoint("LOADPCT", &point)
						case "BCHARGE":
							d64, _ := strconv.ParseFloat(strings.TrimSpace(floatStr[0]), 32)
							d64 = v.chartPageData[h.Name][k].AddValue(d64)
							point := sknlinechart.NewChartDatapoint(float32(d64), theme.ColorGreen, time.Now().Format(time.RFC1123))
							chart.ApplyDataPoint("BCHARGE", &point)
						case "TIMELEFT":
							d64, _ := strconv.ParseFloat(strings.TrimSpace(floatStr[0]), 32)
							d64 = v.chartPageData[h.Name][k].AddValue(d64)
							point := sknlinechart.NewChartDatapoint(float32(d64), theme.ColorPurple, time.Now().Format(time.RFC1123))
							chart.ApplyDataPoint("TIMELEFT", &point)
						case "CUMONBATT":
							d64, _ := strconv.ParseFloat(strings.TrimSpace(floatStr[0]), 32)
							d64 = v.chartPageData[h.Name][k].AddValue(d64)
							point := sknlinechart.NewChartDatapoint(float32(d64), theme.ColorOrange, time.Now().Format(time.RFC1123))
							chart.ApplyDataPoint("CUMONBATT", &point)
						case "STATUS":
							if strings.Contains(vv, "ONLINE") {
								h.State = commons.HostStatusOnline
							} else if strings.Contains(vv, "CHARG") {
								h.State = commons.HostStatusCharging
							} else if strings.Contains(vv, "TEST") {
								h.State = commons.PreferencesIcon
							} else if strings.Contains(vv, "ONBATT") {
								h.State = commons.HostStatusOnBattery
							}
						}

					}
					// details page updates
					knowledge <- params

					// ready next cycle
					currentSt = currentSt[:0]
				}
			}
		}
		// cleanup data syncs
		commons.DebugLog("ViewProvider::HandleUpdatesForMonitorPage[", h.Name, "] ENDED")
	}(host, service, status, events, chart, kChan)
}
