package usecase

import (
	"context"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/ktr03rtk/go-gps-logger/receiver/domain/model"
	"github.com/ktr03rtk/go-gps-logger/receiver/domain/repository"
	"github.com/pkg/errors"
)

type GpsReceiveUsecase interface {
	Execute(context.Context, time.Duration, chan<- *model.Gps) error
}

type gpsReceiveUsecase struct {
	gpsRepository repository.GpsReceiveRepository
	Clock         clock.Clock
}

func NewGpsReceiveUsecase(gr repository.GpsReceiveRepository, clock clock.Clock) GpsReceiveUsecase {
	return &gpsReceiveUsecase{
		gpsRepository: gr,
		Clock:         clock,
	}
}

func (gu *gpsReceiveUsecase) Execute(ctx context.Context, interval time.Duration, ch chan<- *model.Gps) error {
	ticker := gu.Clock.Ticker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			close(ch)

			return nil
		case <-ticker.C:
			g, err := gu.gpsRepository.Receive()
			if err != nil {
				return errors.Wrap(err, "failed to receive GPS data")
			}

			ch <- g
		}
	}
}
