package main

import (
	"context"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"github.com/skoona/ggapcmon/internal/commons"
	"github.com/skoona/ggapcmon/internal/providers"
	"github.com/skoona/ggapcmon/internal/services"
	"github.com/skoona/ggapcmon/internal/ui"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var err error
	commons.ShutdownSignals = make(chan os.Signal, 1)

	ctx, cancelApc := context.WithCancel(context.Background())
	defer cancelApc()

	gui := app.NewWithID("net.skoona.projects.ggapcmon")
	commons.DebugLog("main()::RootURI: ", gui.Storage().RootURI().Path())
	gui.SetIcon(commons.SknSelectThemedResource(commons.AppIcon))

	go func(stopFlag chan os.Signal, a fyne.App) {
		signal.Notify(stopFlag, syscall.SIGINT, syscall.SIGTERM)
		sig := <-stopFlag // wait on ctrl-c
		cancelApc()
		time.Sleep(5 * time.Second)
		err = fmt.Errorf("Shutdown Signal Received: %v", sig.String())
		a.Quit()
	}(commons.ShutdownSignals, gui)

	cfg, err := providers.NewConfig(gui.Preferences())
	if err != nil {
		dialog.ShowError(fmt.Errorf("main()::NewConfig(): %v", err), gui.NewWindow("ggapcmon Configuration Failed"))
		commons.ShutdownSignals <- syscall.SIGINT
		cfg.ResetConfig()
	}

	///*
	//hubHost := entities.NewHubHost("Scotts", "10.100.1.41", "a79c07db-9178-4976-bd10-428aa0d3d159", "10.100.1.183")
	//hubHost := cfg.HubHosts()[0]
	//api := providers.NewHubitatProvider(ctx, hubHost)

	//deviceList := api.DeviceList()
	//commons.DebugLog("DeviceList ", deviceList)

	//deviceInfo := api.DeviceDetailsList()
	//commons.DebugLog("DeviceInfo ", deviceInfo)
	//
	//device := api.DeviceDetailById("3")
	//commons.DebugLog("Device ", device)
	//
	//deviceCapabilities := api.DeviceCapabilitiesById("3")
	//commons.DebugLog("DeviceCapabilities ", deviceCapabilities)
	//
	//deviceHistory := api.DeviceEventHistoryById("7")
	//commons.DebugLog("DeviceHistory ", deviceHistory)

	//ok := api.CreateDeviceEventListener()
	//commons.DebugLog("CreateDeviceEventListener ", ok)
	//if ok {
	//	time.Sleep(1 * time.Minute)
	//	cancelApc()
	//	gui.Quit()
	//	commons.DebugLog("HubHost DeviceDetails ==> ", hubHost.DeviceDetails)
	//}
	//time.Sleep(1 * time.Second)
	//return
	//*/

	service, err := services.NewService(ctx, cfg)
	if err != nil {
		log.Panic("main()::Service startup() failed: ", err.Error())
	}
	defer service.Shutdown()

	vp := ui.NewViewProvider(ctx, cfg, service)
	defer vp.Shutdown()

	vp.ShowMainPage()
	gui.Run()
	commons.DebugLog("main::Shutdown Ended ")
}
