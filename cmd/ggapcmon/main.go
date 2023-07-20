package main

import (
	"context"
	"fmt"
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

	go func(stopFlag chan os.Signal) {
		signal.Notify(stopFlag, syscall.SIGINT, syscall.SIGTERM)
		sig := <-stopFlag // wait on ctrl-c
		cancelApc()
		time.Sleep(5 * time.Second)
		err = fmt.Errorf("Shutdown Signal Received: %v", sig.String())
	}(shutdownSignals)

	service, err := services.NewService(ctx, cfg, logger)
	if err != nil {
		log.Panic("main()::Service startup() failed: ", err.Error())
	}
	defer service.Shutdown()

	vp := ui.NewViewProvider(ctx, cfg, service, logger)
	defer vp.Shutdown()

	go func() {
		logger.Println("main::Shutdown Listener BEGIN ")
	basic:
		for {
			select {
			case <-ctx.Done():
				logger.Println("main::Done() fired:", ctx.Err().Error())
				err = ctx.Err()
				time.Sleep(50 * time.Millisecond)
				gui.Quit()
				break basic

			case msg := <-service.HostMessageChannel(providers.HostPveName).Status:
				for idx, item := range msg {
					logger.Print("{", providers.HostPveName, "}", "(", idx, ")[status] ==> ", item)
				}
			case msg := <-service.HostMessageChannel(providers.HostPveName).Events:
				for idx, item := range msg {
					logger.Print("{", providers.HostPveName, "}", "(", idx, ")[events] ==> ", item)
				}
			case msg := <-service.HostMessageChannel(providers.HostVServName).Status:
				for idx, item := range msg {
					logger.Print("{", providers.HostVServName, "}", "(", idx, ")[status] ==> ", item)
				}
			case msg := <-service.HostMessageChannel(providers.HostVServName).Events:
				for idx, item := range msg {
					logger.Print("{", providers.HostVServName, "}", "(", idx, ")[events] ==> ", item)
				}
			}
		}

		logger.Println("main::Shutdown Listener END ", err.Error())
	}()

	vp.ShowMainPage()
	gui.Run()
	logger.Println("main::Shutdown Ended ")
}
