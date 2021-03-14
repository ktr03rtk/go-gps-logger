package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ktr03rtk/go-gps-logger/pkg/datacommunicator"
)

func main() {
	a := *datacommunicator.NewCommunicator()
	a.Communicate()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		for {
			data := a.Receive()

			log.Printf("%%#v (%#v)", data)
		}
	}()

	<-sig

	a.Close()
}