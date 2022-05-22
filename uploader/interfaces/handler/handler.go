package handler

import (
	"context"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/ktr03rtk/go-gps-logger/uploader/usecase"
)

type LogUploadHandler interface {
	Handle(ctx context.Context, interval time.Duration) error
}

type logUploadHandler struct {
	fileReadUsecase      usecase.FileReadUsecase
	payloadUploadUsecase usecase.PayloadUploadUsecase
	fileDeleteUsecase    usecase.FileDeleteUsecase
	Clock                clock.Clock
}

func NewLogUploadHandler(ru usecase.FileReadUsecase, uu usecase.PayloadUploadUsecase, du usecase.FileDeleteUsecase, clock clock.Clock) LogUploadHandler {
	return &logUploadHandler{
		fileReadUsecase:      ru,
		payloadUploadUsecase: uu,
		fileDeleteUsecase:    du,
		Clock:                clock,
	}
}

func (lh logUploadHandler) Handle(ctx context.Context, interval time.Duration) error {
	ticker := lh.Clock.Ticker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():

			return nil
		case <-ticker.C:
			p, err := lh.fileReadUsecase.Execute()
			if err != nil {
				return err
			}

			filePaths, err := lh.payloadUploadUsecase.Execute(ctx, p)
			if err != nil {
				return err
			}

			if err := lh.fileDeleteUsecase.Execute(filePaths); err != nil {
				return err
			}
		}
	}
}
