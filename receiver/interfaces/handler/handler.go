package handler

import (
	"context"
	"time"

	"github.com/ktr03rtk/go-gps-logger/receiver/domain/model"
	"github.com/ktr03rtk/go-gps-logger/receiver/usecase"
	"golang.org/x/sync/errgroup"
)

type GpsLogHandler interface {
	Handle(ctx context.Context, interval time.Duration) error
}

type gpsLogHandler struct {
	gpsReceiveUsecase usecase.GpsReceiveUsecase
	gpsWriteUsecase   usecase.GpsWriteUsecase
}

func NewGpsLogHandler(ru usecase.GpsReceiveUsecase, wu usecase.GpsWriteUsecase) GpsLogHandler {
	return &gpsLogHandler{
		gpsReceiveUsecase: ru,
		gpsWriteUsecase:   wu,
	}
}

func (gh gpsLogHandler) Handle(ctxParent context.Context, interval time.Duration) error {
	eg, ctx := errgroup.WithContext(ctxParent)
	ch := make(chan *model.Gps, 10)

	eg.Go(func() error { return gh.gpsReceiveUsecase.Execute(ctx, interval, ch) })
	eg.Go(func() error { return gh.gpsWriteUsecase.Execute(ctx, ch) })

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}
