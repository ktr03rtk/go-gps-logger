package usecase

import (
	"context"

	"github.com/ktr03rtk/go-gps-logger/receiver/domain/model"
	"github.com/ktr03rtk/go-gps-logger/receiver/domain/repository"
	"github.com/pkg/errors"
)

type GpsWriteUsecase interface {
	Execute(context.Context, <-chan *model.Gps) error
}

type gpsWriteUsecase struct {
	gpsRepository repository.GpsWriteRepository
}

func NewGpsWriteUsecase(gr repository.GpsWriteRepository) GpsWriteUsecase {
	return &gpsWriteUsecase{gpsRepository: gr}
}

func (gu *gpsWriteUsecase) Execute(ctx context.Context, ch <-chan *model.Gps) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case gps := <-ch:
			if err := gu.gpsRepository.Write(gps); err != nil {
				if err != nil {
					return errors.Wrap(err, "failed to write GPS data")
				}
			}
		}
	}
}
