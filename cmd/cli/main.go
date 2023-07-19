package main

import (
	"context"
	"fmt"
	"github.com/skoona/ggapcmon/internal/entities"
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

	go func(stopFlag chan os.Signal) {
		signal.Notify(stopFlag, syscall.SIGINT, syscall.SIGTERM)
		sig := <-stopFlag // wait on ctrl-c
		cancelApc()
		time.Sleep(5 * time.Second)
		err = fmt.Errorf("Shutdown Signal Received: %v", sig.String())
	}(shutdownSignals)

	// get from fyne.Preferences
	hosts := []entities.ApcHost{
		{IpAddress: HostVServ, Name: HostVServName, SecondsPerSample: 33},
		{IpAddress: HostPve, Name: HostPveName, SecondsPerSample: 37},
	}

	service, err := services.NewService(ctx, hosts, logger)
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

		case msg := <-service.HostMessageChannel(hosts[0].Name):
			for idx, item := range msg {
				parts := strings.SplitN(item, ": ", 2)
				logger.Print("{", hosts[0].Name, "}", "(", idx, ")[", parts[0], "] ==> ", parts[1])
			}
		case msg := <-service.HostMessageChannel(hosts[1].Name):
			for idx, item := range msg {
				parts := strings.SplitN(item, ": ", 2)
				logger.Print("{", hosts[1].Name, "}", "(", idx, ")[", parts[0], "] ==> ", parts[1])
			}
		}
	}

	logger.Println("main::Shutdown", err.Error())
}
