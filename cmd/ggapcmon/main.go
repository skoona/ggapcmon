package main

import (
	"context"
	"fmt"
	"fyne.io/fyne/v2/app"
	"github.com/skoona/ggapcmon/internal/services"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	HostVServ     = "10.100.1.3:3551"
	HostVServName = "VServ"
	HostPve       = "10.100.1.4:3551"
	HostPveName   = "PVE"
)

func main() {
	var err error
	logger := log.New(os.Stdout, "[DEBUG] ", log.Lmicroseconds|log.Lshortfile)
	shutdownSignals := make(chan os.Signal, 1)

	ctx, cancelApc := context.WithCancel(context.Background())
	defer cancelApc()

	gui := app.NewWithID("net.skoona.project.ggapcmon")
	logger.Print("main()::RootURI: ", gui.Storage().RootURI().Path())

	cfg, err := services.NewConfig(gui.Preferences(), logger)
	if err != nil {
		log.Panic("main()::NewConfig() failed: ", err.Error())
	}

	go func(stopFlag chan os.Signal) {
		signal.Notify(stopFlag, syscall.SIGINT, syscall.SIGTERM)
		sig := <-stopFlag // wait on ctrl-c
		cancelApc()
		time.Sleep(5 * time.Second)
		err = fmt.Errorf("Shutdown Signal Received: %v", sig.String())
	}(shutdownSignals)

	service, err := services.NewService(ctx, cfg.Hosts(), logger)
	if err != nil {
		log.Panic("main()::Service startup() failed: ", err.Error())
	}
	defer service.Shutdown()

basic:
	for {
		select {
		case <-ctx.Done():
			logger.Println("main::Done() fired:", ctx.Err().Error())
			err = ctx.Err()
			break basic

		case msg := <-service.HostMessageChannel(services.HostPveName):
			for idx, item := range msg {
				parts := strings.SplitN(item, ": ", 2)
				logger.Print("{", services.HostPveName, "}", "(", idx, ")[", parts[0], "] ==> ", parts[1])
			}
		case msg := <-service.HostMessageChannel(services.HostVServName):
			for idx, item := range msg {
				parts := strings.SplitN(item, ": ", 2)
				logger.Print("{", services.HostVServName, "}", "(", idx, ")[", parts[0], "] ==> ", parts[1])
			}
		}
	}

	logger.Println("main::Shutdown", err.Error())
}
