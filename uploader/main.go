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
	"github.com/ktr03rtk/go-gps-logger/uploader/infrastructure/adapter"
	"github.com/ktr03rtk/go-gps-logger/uploader/infrastructure/persistence"
	"github.com/ktr03rtk/go-gps-logger/uploader/interfaces/handler"
	"github.com/ktr03rtk/go-gps-logger/uploader/usecase"
)

var (
	sourceDir      string
	uploadInterval time.Duration
)

func init() {
	if err := getEnv(); err != nil {
		log.Fatal(err)
	}
}

func getEnv() error {
	d, ok := os.LookupEnv("SOURCE_DIRECTORY")
	if !ok {
		return errors.New("env SOURCE_DIRECTORY is not found")
	}

	sourceDir = d

	r, ok := os.LookupEnv("UPLOAD_INTERVAL_SECONDS")
	if !ok {
		return errors.New("env UPLOAD_INTERVAL_SECONDS is not found")
	}

	i, err := strconv.Atoi(r)
	if err != nil {
		return errors.New("env UPLOAD_INTERVAL_SECONDS is not integer")
	}

	uploadInterval = time.Duration(i) * time.Second

	return nil
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	fr := persistence.NewFilePersistence(sourceDir)
	pr, err := adapter.NewMqttAdapter(ctx)
	if err != nil {
		log.Fatal(err)
	}

	ru := usecase.NewFileReadUsecase(fr)
	uu := usecase.NewPayloadUploadUsecase(pr)
	du := usecase.NewFileDeleteUsecase(fr)
	h := handler.NewLogUploadHandler(ru, uu, du, clock.New())

	if err := h.Handle(ctx, uploadInterval); err != nil {
		log.Fatal(err)
	}
}
