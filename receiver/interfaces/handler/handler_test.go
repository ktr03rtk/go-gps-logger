package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/golang/mock/gomock"
	"github.com/ktr03rtk/go-gps-logger/receiver/domain/model"
	"github.com/ktr03rtk/go-gps-logger/receiver/mock"
	"github.com/ktr03rtk/go-gps-logger/receiver/usecase"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func TestGpsReceiveUsecaseExecute(t *testing.T) {
	date := time.Date(2022, 5, 3, 0, 9, 0, 0, time.Local)
	duration := 1000 * time.Millisecond

	tests := []struct {
		name             string
		receiveOutput    *model.Gps
		receiveErr       error
		writeErr         error
		expectedErr      error
		receiveCallTimes int
		writeCallTimes   int
	}{
		{
			"normal case",
			&model.Gps{Timestamp: date, Mode: 3, Lat: 30.11, Lon: 130.11, Alt: 30.11, Speed: 30.11},
			nil,
			nil,
			nil,
			1,
			1,
		},
		{
			"receive error case",
			nil,
			errors.New("receive repository error"),
			nil,
			errors.New("failed to receive GPS data"),
			1,
			0,
		},
		{
			"write error case",
			&model.Gps{Timestamp: date, Mode: 3, Lat: 30.11, Lon: 130.11, Alt: 30.11, Speed: 30.11},
			nil,
			errors.New("write repository error"),
			errors.New("failed to write GPS data"),
			1,
			1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rr := mock.NewMockGpsReceiveRepository(ctrl)
			wr := mock.NewMockGpsWriteRepository(ctrl)
			mock := clock.NewMock()
			ru := usecase.NewGpsReceiveUsecase(rr, mock)
			wu := usecase.NewGpsWriteUsecase(wr)
			h := NewGpsLogHandler(ru, wu)

			rr.EXPECT().Receive().Return(tt.receiveOutput, tt.receiveErr).Times(tt.receiveCallTimes)
			wr.EXPECT().Write(tt.receiveOutput).Return(tt.writeErr).Times(tt.writeCallTimes)

			ctx, cancel := context.WithCancel(context.Background())
			eg, ctx := errgroup.WithContext(ctx)
			eg.Go(func() error { return h.Handle(ctx, duration) })

			gosched()
			mock.Add(duration)

			cancel()

			if err := eg.Wait(); err != nil {
				if tt.expectedErr != nil {
					assert.Contains(t, err.Error(), tt.expectedErr.Error())
				} else {
					t.Fatalf("error is not expected but received: %v", err)
				}
			} else {
				assert.Exactly(t, tt.expectedErr, nil, "error is expected but received nil")
			}
		})
	}
}

// Sleep momentarily so that other goroutines can process.
func gosched() { time.Sleep(1 * time.Millisecond) }
