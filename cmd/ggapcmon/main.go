package main

import (
	"context"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
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
	shutdownSignals := make(chan os.Signal, 1)

	ctx, cancelApc := context.WithCancel(context.Background())
	defer cancelApc()

	gui := app.NewWithID("net.skoona.project.ggapcmon")
	logger.Print("main()::RootURI: ", gui.Storage().RootURI().Path())

	cfg, err := providers.NewConfig(gui.Preferences(), logger)
	if err != nil {
		log.Panic("main()::NewConfig() failed: ", err.Error())
	}
	//cfg.ResetConfig()
	//logger.Println("HostKeys: ", cfg.HostKeys())

	go func(stopFlag chan os.Signal, a fyne.App) {
		signal.Notify(stopFlag, syscall.SIGINT, syscall.SIGTERM)
		sig := <-stopFlag // wait on ctrl-c
		cancelApc()
		time.Sleep(5 * time.Second)
		err = fmt.Errorf("Shutdown Signal Received: %v", sig.String())
		a.Quit()
	}(shutdownSignals, gui)

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
