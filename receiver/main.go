package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/ktr03rtk/go-gps-logger/receiver/infrastructure/adapter"
	"github.com/ktr03rtk/go-gps-logger/receiver/infrastructure/persistence"
	"github.com/ktr03rtk/go-gps-logger/receiver/interfaces/handler"
	"github.com/ktr03rtk/go-gps-logger/receiver/usecase"
)

var (
	distDir  string
	interval time.Duration
)

func init() {
	if err := getEnv(); err != nil {
		log.Fatal(err)
	}
}

func getEnv() error {
	d, ok := os.LookupEnv("DIST_DIRECTORY")
	if !ok {
		return errors.New("env DIST_DIRECTORY is not found")
	}

	distDir = d

	r, ok := os.LookupEnv("INTERVAL_SECONDS")
	if !ok {
		return errors.New("env INTERVAL_SECONDS is not found")
	}

	i, err := strconv.Atoi(r)
	if err != nil {
		return errors.New("env INTERVAL_SECONDS is not integer")
	}

	interval = time.Duration(i) * time.Second

	return nil
}

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
