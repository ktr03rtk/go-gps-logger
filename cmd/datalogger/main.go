package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ktr03rtk/go-gps-logger/pkg/datacommunicator"
	"github.com/ktr03rtk/go-gps-logger/pkg/dataconverter"
	"github.com/ktr03rtk/go-gps-logger/pkg/datawriter"
)

func main() {
	com := *datacommunicator.NewCommunicator()
	com.Communicate()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		for {
			data := com.Receive()

			converted_data := dataconverter.Convert(data)

			datawriter.Write(converted_data)

		}
	}()

	<-sig

	com.Close()
}
