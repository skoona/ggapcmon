package main

import (
	"context"
	"fmt"
	"github.com/skoona/ggapcmon/internal/services"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	HostVserv = "10.100.1.3:3551"
	HostPve   = "10.100.1.4:3551"
)

func main() {
	systemSignalChannel := make(chan os.Signal, 1)
	msgs := make(chan string, 100)
	var err error

	ctx, cancelApc := context.WithCancel(context.Background())
	defer cancelApc()

	logger := log.New(os.Stdout, "[DEBUG] ", log.Lmicroseconds|log.Lshortfile)

	go func(stopFlag chan os.Signal) {
		signal.Notify(stopFlag, syscall.SIGINT, syscall.SIGTERM)
		sig := <-stopFlag // wait on ctrl-c
		cancelApc()
		time.Sleep(5 * time.Second)
		err = fmt.Errorf("Shutdown Signal Received: %v", sig.String())
	}(systemSignalChannel)

	apc := services.NewServer(ctx, "VServ", HostVserv, 5)
	defer apc.PeriodicUpdateStop()
	defer apc.End()
	apc.Begin()
	err = apc.PeriodicUpdateStart(msgs)

basic:
	for {
		select {
		case <-ctx.Done():
			logger.Println("main::Done() fired:", ctx.Err().Error())
			err = ctx.Err()
			apc.PeriodicUpdateStop()
			apc.End()
			break basic

		case msg := <-msgs:
			parts := strings.SplitN(msg, ": ", 2)
			logger.Print("[", parts[0], "] ==> ", parts[1])
		}
	}
	logger.Println("Shutdown:", err.Error())
}
