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
	logger := log.New(os.Stdout, "[DEBUG] ", log.Lmicroseconds|log.Lshortfile)
	commons.ShutdownSignals = make(chan os.Signal, 1)

	ctx, cancelApc := context.WithCancel(context.Background())
	defer cancelApc()

	gui := app.NewWithID("net.skoona.project.ggapcmon")
	logger.Print("main()::RootURI: ", gui.Storage().RootURI().Path())
	gui.SetIcon(commons.SknSelectThemedResource(commons.AppIcon))

	go func(stopFlag chan os.Signal, a fyne.App) {
		signal.Notify(stopFlag, syscall.SIGINT, syscall.SIGTERM)
		sig := <-stopFlag // wait on ctrl-c
		cancelApc()
		time.Sleep(5 * time.Second)
		err = fmt.Errorf("Shutdown Signal Received: %v", sig.String())
		a.Quit()
	}(commons.ShutdownSignals, gui)

	cfg, err := providers.NewConfig(gui.Preferences(), logger)
	if err != nil {
		dialog.ShowError(fmt.Errorf("main()::NewConfig(): %v", err), gui.NewWindow("ggapcmon Configuration Failed"))
		commons.ShutdownSignals <- syscall.SIGINT
		cfg.ResetConfig()
	}

	service, err := services.NewService(ctx, cfg, logger)
	if err != nil {
		log.Panic("main()::Service startup() failed: ", err.Error())
	}
	defer service.Shutdown()

	vp := ui.NewViewProvider(ctx, cfg, service, logger)
	defer vp.Shutdown()

	vp.ShowMainPage()
	gui.Run()
	logger.Println("main::Shutdown Ended ")
}