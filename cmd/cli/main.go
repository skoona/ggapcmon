package main

import (
	"context"
	"fmt"
	"github.com/skoona/ggapcmon/internal/providers"
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
	msgs := make(chan []string, 16)
	defer close(msgs)
	msgsB := make(chan []string, 16)
	defer close(msgsB)

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

	apc, err := providers.NewAPCProvider(ctx, "VServ", HostVserv, 6, msgs)
	if err != nil {
		log.Panic("main()::Begin(A) failed: ", err.Error())
	}
	defer apc.Shutdown()

	apcB, err := providers.NewAPCProvider(ctx, "PVE", HostPve, 7, msgsB)
	if err != nil {
		log.Panic("main()::Begin(B) failed: ", err.Error())
	}
	defer apcB.Shutdown()

basic:
	for {
		select {
		case <-ctx.Done():
			logger.Println("main::Done() fired:", ctx.Err().Error())
			err = ctx.Err()
			break basic

		case msg := <-msgsB:
			for idx, item := range msg {
				parts := strings.SplitN(item, ": ", 2)
				logger.Print("{", apc.Name(), "}", "(", idx, ")[", parts[0], "] ==> ", parts[1])
			}
		case msg := <-msgs:
			for idx, item := range msg {
				parts := strings.SplitN(item, ": ", 2)
				logger.Print("{", apcB.Name(), "}", "(", idx, ")[", parts[0], "] ==> ", parts[1])
			}
		}
	}

	logger.Println("Shutdown:", err.Error())
}
