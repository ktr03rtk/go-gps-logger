package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/ktr03rtk/go-gps-logger/receiver/infrastructure/adapter"
	"github.com/ktr03rtk/go-gps-logger/receiver/infrastructure/persistence"
	"github.com/ktr03rtk/go-gps-logger/receiver/interfaces/handler"
	"github.com/ktr03rtk/go-gps-logger/receiver/usecase"
)

const (
	distDir  = "/var/data/gps/raw"
	interval = 10 * time.Second
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	rr, err := adapter.NewGpsAdapter(ctx)
	if err != nil {
		log.Fatal(err)
	}

	wr, err := persistence.NewGpsPersistence(distDir)
	if err != nil {
		log.Fatal(err)
	}

	ru := usecase.NewGpsReceiveUsecase(rr, clock.New())
	wu := usecase.NewGpsWriteUsecase(wr)
	h := handler.NewGpsLogHandler(ru, wu)

	if err := h.Handle(ctx, interval); err != nil {
		log.Fatal(err)
	}
}
