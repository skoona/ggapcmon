package main

import (
	"context"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"github.com/skoona/ggapcmon/internal/adapters/handlers/ui"
	"github.com/skoona/ggapcmon/internal/adapters/repository"
	"github.com/skoona/ggapcmon/internal/commons"
	"github.com/skoona/ggapcmon/internal/core/services"
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
		err = fmt.Errorf("Close Signal Received: %v", sig.String())
		cancelApc()
		time.Sleep(5 * time.Second)
		a.Quit()
	}(commons.ShutdownSignals, gui)

	cfg, err := repository.NewConfig(gui.Preferences())
	if err != nil {
		dialog.ShowError(fmt.Errorf("main()::NewConfig(): %v", err), gui.NewWindow("ggapcmon Configuration Failed"))
		commons.ShutdownSignals <- syscall.SIGINT
		cfg.ResetConfig()
	}
	//cfg.ResetConfig()

	service, err := services.NewService(ctx, cfg)
	if err != nil {
		log.Panic("main()::Service startup() failed: ", err.Error())
	}
	defer service.Close()

	vp := ui.NewViewProvider(ctx, cfg, service)
	defer vp.Close()

	vp.ShowMainPage()
	gui.Run()
	commons.DebugLog("main::Close Ended ")
}
