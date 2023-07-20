package ui

import (
	"context"
	"fyne.io/fyne/v2"
	"github.com/skoona/ggapcmon/internal/entities"
	"github.com/skoona/ggapcmon/internal/interfaces"
	"log"
)

type viewProvider struct {
	ctx          context.Context
	service      interfaces.Service
	mainWindow   fyne.Window
	prefsWindow  fyne.Window
	log          *log.Logger
	cfg          interfaces.Configuration
	prfHostKeys  []string
	prfHost      entities.ApcHost
	prfAddAction bool
	prfDelAction bool
}

var (
	_ interfaces.ViewProvider = (*viewProvider)(nil)
	_ interfaces.Provider     = (*viewProvider)(nil)
)

func NewViewProvider(ctx context.Context, cfg interfaces.Configuration, service interfaces.Service, log *log.Logger) interfaces.ViewProvider {
	view := &viewProvider{
		ctx:         ctx,
		log:         log,
		cfg:         cfg,
		service:     service,
		mainWindow:  fyne.CurrentApp().NewWindow("ggAPC Monitor"),
		prefsWindow: fyne.CurrentApp().NewWindow("Preferences"),
	}
	view.mainWindow.Resize(fyne.NewSize(1024, 756))
	view.mainWindow.SetCloseIntercept(func() { view.mainWindow.Hide() })
	view.mainWindow.SetMaster()

	view.prefsWindow.Resize(fyne.NewSize(632, 512))
	view.prefsWindow.SetCloseIntercept(func() { view.prefsWindow.Hide() })

	view.SknTrayMenu()
	view.SknMenus()

	return view
}
func (v *viewProvider) ShowMainPage() {
	v.mainWindow.SetContent(v.MonitorPage())
	v.mainWindow.Show()
}

func (v *viewProvider) ShowPrefsPage() {
	v.prefsWindow.SetContent(v.PrefsPage())
	v.prefsWindow.Show()
}

// Shutdown closes all go routine
func (a *viewProvider) Shutdown() {
	a.log.Println("ViewProvider::Shutdown() called.")
}

// prefsAddAction closes all go routine
func (v *viewProvider) prefsAddAction() {
	v.prfAddAction = true
	v.cfg.AddHost(v.prfHost)
	v.ShowPrefsPage()
	v.log.Println("ViewProvider::prefsAddAction() called.")
}

// prefsDelAction closes all go routine
func (v *viewProvider) prefsDelAction() {
	v.prfDelAction = true
	v.cfg.Remove(v.prfHost.Name)
	v.ShowPrefsPage()
	v.log.Println("ViewProvider::prefsDelAction() called.")
}
