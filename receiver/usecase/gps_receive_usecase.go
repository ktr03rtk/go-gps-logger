package usecase

import (
	"context"
	"time"

	"github.com/ktr03rtk/go-gps-logger/receiver/domain/model"
	"github.com/ktr03rtk/go-gps-logger/receiver/domain/repository"
	"github.com/pkg/errors"
)

type GpsReceiveUsecase interface {
	Execute(context.Context, time.Duration, chan<- *model.Gps) error
}

type gpsReceiveUsecase struct {
	gpsRepository repository.GpsReceiveRepository
}

func NewGpsReceiveUsecase(gr repository.GpsReceiveRepository) GpsReceiveUsecase {
	return &gpsReceiveUsecase{gpsRepository: gr}
}

func (u *gpsReceiveUsecase) Execute(ctx context.Context, interval time.Duration, ch chan<- *model.Gps) error {
	for {
		select {
		case <-ctx.Done():
			close(ch)

			return nil
		case <-time.Tick(interval):
			g, err := u.gpsRepository.Receive()
			if err != nil {
				return errors.Wrap(err, "failed to receive GPS data")
			}

			ch <- g
		}
	}
}
