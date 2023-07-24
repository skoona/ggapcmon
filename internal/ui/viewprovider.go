package ui

import (
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/skoona/ggapcmon/internal/commons"
	"github.com/skoona/ggapcmon/internal/entities"
	"github.com/skoona/ggapcmon/internal/interfaces"
)

// viewProvider control structure for view management
type viewProvider struct {
	ctx             context.Context
	cfg             interfaces.Configuration
	service         interfaces.Service
	mainWindow      fyne.Window
	prefsWindow     fyne.Window
	prfStatusLine   *widget.Label
	chartKeys       []string
	chartPageData   map[string]map[string]interfaces.GraphPointSmoothing
	bondedUpsStatus map[string]*entities.UpsStatusValueBindings
	prfHostKeys     []string
	prfHost         *entities.ApcHost
}

// compiler helpers to insure interfaces requirements are meet
var (
	_ interfaces.ViewProvider = (*viewProvider)(nil)
	_ interfaces.Provider     = (*viewProvider)(nil)
)

// NewViewProvider manages all UI views and implements the ViewProvider Interface
func NewViewProvider(ctx context.Context, cfg interfaces.Configuration, service interfaces.Service) interfaces.ViewProvider {
	hk := cfg.HostKeys()
	h := cfg.HostByName(hk[0])
	stLine := widget.NewLabel("click entry in table to edit, or click add to add.")
	//stLine.Wrapping = fyne.TextWrapWord -- causes pref page to break
	view := &viewProvider{
		ctx:             ctx,
		cfg:             cfg,
		service:         service,
		mainWindow:      fyne.CurrentApp().NewWindow("ggAPC Monitor"),
		prefsWindow:     fyne.CurrentApp().NewWindow("Preferences"),
		prfHost:         h,
		prfHostKeys:     hk,
		prfStatusLine:   stLine,
		chartPageData:   map[string]map[string]interfaces.GraphPointSmoothing{}, // [host][chartkey]struct
		chartKeys:       []string{"LINEV", "LOADPCT", "BCHARGE", "CUMONBATT", "TIMELEFT"},
		bondedUpsStatus: map[string]*entities.UpsStatusValueBindings{},
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

// ShowMainPage display the primary application page
func (v *viewProvider) ShowMainPage() {
	v.prfHostKeys = v.cfg.HostKeys()
	v.mainWindow.SetContent(v.MonitorPage())
	v.mainWindow.Show()
}

// ShowPrefsPage displays teh settings por preferences page
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
