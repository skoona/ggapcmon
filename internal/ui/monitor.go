package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/skoona/ggapcmon/internal/commons"

	"github.com/skoona/sknlinechart"
	"syscall"
)

// monitor page

func (v *viewProvider) MonitorPage() *fyne.Container {
	hostTabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Known UPSs", theme.ComputerIcon(), v.OverviewPage()),
		container.NewTabItemWithIcon("Glossary", theme.InfoIcon(), v.GlossaryPage()),
	)

	for _, host := range v.cfg.Hosts() {
		if host.Enabled {
			status := widget.NewMultiLineEntry()
			status.SetPlaceHolder("Status Page")
			status.TextStyle = fyne.TextStyle{Monospace: true}

			events := widget.NewMultiLineEntry()
			events.SetPlaceHolder("Events Page")
			events.TextStyle = fyne.TextStyle{Monospace: true}

			// for chart page updates
			data := map[string][]*sknlinechart.ChartDatapoint{}
			chart, err := sknlinechart.NewLineChart(
				host.Name,
				"",
				&data,
			)
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
					container.NewTabItemWithIcon("Detailed", theme.VisibilityIcon(), v.DetailPage(knowledge)),
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
