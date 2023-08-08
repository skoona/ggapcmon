package ui

import (
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/skoona/ggapcmon/internal/commons"
	"github.com/skoona/ggapcmon/internal/core/domain"
	"github.com/skoona/ggapcmon/internal/core/ports"
)

type ViewProvider interface {
	ShowPrefsPage()
	ShowMainPage()
	ports.Provider
}

// viewProvider control structure for view management
type viewProvider struct {
	ctx             context.Context
	cfg             ports.Configuration
	service         ports.Service
	mainWindow      fyne.Window
	prefsWindow     fyne.Window
	overviewTable   *widget.Table
	prfStatusLine   *widget.Label
	chartKeys       []string
	chartPageData   map[string]map[string]ports.GraphPointSmoothing
	bondedUpsStatus map[string]*domain.UpsStatusValueBindings
	prfHostKeys     []string
	prfHost         *domain.ApcHost
}

// compiler helpers to insure ports requirements are meet
var (
	_ ViewProvider   = (*viewProvider)(nil)
	_ ports.Provider = (*viewProvider)(nil)
)

// NewViewProvider manages all UI views and implements the ViewProvider Interface
func NewViewProvider(ctx context.Context, cfg ports.Configuration, service ports.Service) ViewProvider {
	hk := cfg.HostKeys()
	h := cfg.HostById(hk[0])
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
		chartPageData:   map[string]map[string]ports.GraphPointSmoothing{}, // [host][chartkey]struct
		chartKeys:       []string{"LINEV", "LOADPCT", "BCHARGE", "CUMONBATT", "TIMELEFT"},
		bondedUpsStatus: map[string]*domain.UpsStatusValueBindings{},
	}
	view.mainWindow.Resize(fyne.NewSize(960, 496))
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
	v.prfHost = v.cfg.HostById(v.prfHostKeys[0])
	v.prefsWindow.SetContent(v.PreferencesPage())
	v.prefsWindow.Show()
}

// Close closes all go routine
func (v *viewProvider) Close() {
	commons.DebugLog("ViewProvider::Close() called.")
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
