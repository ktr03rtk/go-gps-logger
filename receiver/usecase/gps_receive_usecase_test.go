package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/golang/mock/gomock"
	"github.com/ktr03rtk/go-gps-logger/receiver/domain/model"
	"github.com/ktr03rtk/go-gps-logger/receiver/mock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func TestGpsReceiveUsecaseExecute(t *testing.T) {
	t.Parallel()

	date := time.Date(2022, 5, 3, 0, 9, 0, 0, time.Local)
	duration := 1000 * time.Millisecond

	tests := []struct {
		name              string
		returnErr         error
		expectedOutput    *model.Gps
		expectedErr       error
		expectedCallTimes int
	}{
		{
			"normal case",
			nil,
			&model.Gps{Timestamp: date, Mode: 3, Lat: 30.11, Lon: 130.11, Alt: 30.11, Speed: 30.11},
			nil,
			2,
		},
		{
			"error case",
			errors.New("error occurred"),
			nil,
			errors.New("failed to receive GPS data"),
			1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repository := mock.NewMockGpsReceiveRepository(ctrl)
			mock := clock.NewMock()
			usecase := NewGpsReceiveUsecase(repository, mock)

			context, cancel := context.WithCancel(context.Background())
			eg, ctx := errgroup.WithContext(context)
			ch := make(chan *model.Gps)

			repository.EXPECT().Receive().Return(tt.expectedOutput, tt.returnErr).Times(tt.expectedCallTimes)
			eg.Go(func() error { return usecase.Execute(ctx, duration, ch) })

			for i := 0; i < tt.expectedCallTimes; i++ {
				gosched()
				mock.Add(duration)

				if tt.expectedErr == nil {
					output := <-ch
					assert.Exactly(t, tt.expectedOutput, output)
				}
			}

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

func TestGpsReceiveUsecaseExecuteContexCancelled(t *testing.T) {
	t.Parallel()

	duration := 1000 * time.Millisecond

	tests := []struct {
		name              string
		returnErr         error
		expectedOutput    *model.Gps
		expectedErr       error
		expectedCallTimes int
	}{
		{
			"context canceled case",
			nil,
			nil,
			nil,
			0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repository := mock.NewMockGpsReceiveRepository(ctrl)
			mock := clock.NewMock()
			usecase := NewGpsReceiveUsecase(repository, mock)

			context, cancel := context.WithCancel(context.Background())
			eg, ctx := errgroup.WithContext(context)
			ch := make(chan *model.Gps)

			repository.EXPECT().Receive().Return(tt.expectedOutput, tt.returnErr).Times(tt.expectedCallTimes)
			eg.Go(func() error { return usecase.Execute(ctx, duration, ch) })

			cancel()
			mock.Add(duration)
			gosched()
			_, ok := <-ch

			err := eg.Wait()
			assert.Nil(t, err, "error is not expected but received: %v", err)
			assert.False(t, ok, "channel is expected to be closed")
		})
	}
}

// Sleep momentarily so that other goroutines can process.
func gosched() { time.Sleep(1 * time.Millisecond) }
