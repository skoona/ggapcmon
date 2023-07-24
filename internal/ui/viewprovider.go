package ui

import (
	"context"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/skoona/ggapcmon/internal/commons"
	"github.com/skoona/ggapcmon/internal/entities"
	"github.com/skoona/ggapcmon/internal/interfaces"
	"github.com/skoona/sknlinechart"
	"strconv"
	"strings"
	"time"
)

type viewProvider struct {
	ctx           context.Context
	service       interfaces.Service
	mainWindow    fyne.Window
	prefsWindow   fyne.Window
	prfStatusLine *widget.Label
	cfg           interfaces.Configuration
	chartPageData map[string]map[string]interfaces.GraphPointSmoothing
	chartKeys     []string
	prfHostKeys   []string
	prfHost       *entities.ApcHost
}

var (
	_ interfaces.ViewProvider = (*viewProvider)(nil)
	_ interfaces.Provider     = (*viewProvider)(nil)
)

func NewViewProvider(ctx context.Context, cfg interfaces.Configuration, service interfaces.Service) interfaces.ViewProvider {
	hk := cfg.HostKeys()
	h := cfg.HostByName(hk[0])
	stLine := widget.NewLabel("click entry in table to edit, or click add to add.")
	//stLine.Wrapping = fyne.TextWrapWord -- causes pref page to break
	view := &viewProvider{
		ctx:           ctx,
		cfg:           cfg,
		service:       service,
		mainWindow:    fyne.CurrentApp().NewWindow("ggAPC Monitor"),
		prefsWindow:   fyne.CurrentApp().NewWindow("Preferences"),
		prfHost:       h,
		prfHostKeys:   hk,
		prfStatusLine: stLine,
		chartPageData: map[string]map[string]interfaces.GraphPointSmoothing{}, // [host][chartkey]struct
		chartKeys:     []string{"LINEV", "LOADPCT", "BCHARGE", "CUMONBATT", "TIMELEFT"},
	}
	view.mainWindow.Resize(fyne.NewSize(632, 432))
	view.mainWindow.SetCloseIntercept(func() { view.mainWindow.Hide() })
	view.mainWindow.SetMaster()
	view.mainWindow.SetIcon(commons.SknSelectThemedResource(commons.AppIcon))

	view.prefsWindow.Resize(fyne.NewSize(632, 572))
	view.prefsWindow.SetCloseIntercept(func() { view.prefsWindow.Hide() })
	view.mainWindow.SetIcon(commons.SknSelectThemedResource(commons.PreferencesIcon))

	view.SknTrayMenu()
	view.SknMenus()

	return view
}
func (v *viewProvider) ShowMainPage() {
	v.prfHostKeys = v.cfg.HostKeys()
	v.mainWindow.SetContent(v.MonitorPage())
	v.mainWindow.Show()
}
func (v *viewProvider) ShowPrefsPage() {
	v.prfHostKeys = v.cfg.HostKeys()
	v.prfHost = v.cfg.HostByName(v.prfHostKeys[0])
	v.prefsWindow.SetContent(v.PreferencesPage())
	v.prefsWindow.Show()
}

// Shutdown closes all go routine
func (v *viewProvider) Shutdown() {
	commons.DebugLog("ViewProvider::Shutdown() called.")
}

// verifyHostConnection attempts to connect to selected host
func (v *viewProvider) verifyHostConnection() error {
	err := v.cfg.VerifyHostConnection(v.prfHost)
	if err == nil {
		v.prfStatusLine.SetText("connection to " + v.prfHost.Name + " was successful")
	} else {
		v.prfStatusLine.SetText(v.prfHost.Name + " connect was not successful: " + err.Error())
	}
	return err
}

// prefsAddAction adds or replaces the host in the form
func (v *viewProvider) prefsAddAction() {
	v.prfHostKeys = v.cfg.HostKeys()
	v.ShowPrefsPage()
	v.prfStatusLine.SetText("Host " + v.prfHost.Name + " was added")
}

// prefsDelAction deletes the select host
func (v *viewProvider) prefsDelAction() {
	n := v.prfHost.Name
	v.cfg.Remove(v.prfHost.Name)
	v.prfHostKeys = v.cfg.HostKeys()
	v.prfHost = v.cfg.HostByName(v.prfHostKeys[0])
	v.ShowPrefsPage()
	v.prfStatusLine.SetText("Host " + n + " was removed")
}

// handleUpdatesForMonitorPage does exactly that
func (v *viewProvider) handleUpdatesForMonitorPage(host *entities.ApcHost, service interfaces.Service, status *widget.Entry, events *widget.Entry, chart sknlinechart.LineChart, kChan chan map[string]string) {
	go func(h *entities.ApcHost, svc interfaces.Service, st *widget.Entry, ev *widget.Entry, chart sknlinechart.LineChart, knowledge chan map[string]string) {
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
		commons.DebugLog("ViewProvider::HandleUpdatesForMonitorPage[", h.Name, "] ENDED")
	}(host, service, status, events, chart, kChan)
}
