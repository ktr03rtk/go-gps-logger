package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/ktr03rtk/go-gps-logger/receiver/domain/model"
	"github.com/ktr03rtk/go-gps-logger/receiver/mock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func TestGpsWriteUsecaseExecute(t *testing.T) {
	date := time.Date(2022, 5, 3, 0, 9, 0, 0, time.Local)

	tests := []struct {
		name              string
		inputGps          *model.Gps
		returnErr         error
		expectedErr       error
		expectedCallTimes int
	}{
		{
			"normal case",
			&model.Gps{Timestamp: date, Mode: 3, Lat: 30.11, Lon: 130.11, Alt: 30.11, Speed: 30.11},
			nil,
			nil,
			2,
		},
		{
			"error case",
			nil,
			errors.New("error occurred"),
			errors.New("failed to write GPS data"),
			1,
		},
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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repository := mock.NewMockGpsWriteRepository(ctrl)
			usecase := NewGpsWriteUsecase(repository)

			context, cancel := context.WithCancel(context.Background())
			eg, ctx := errgroup.WithContext(context)
			ch := make(chan *model.Gps)

			repository.EXPECT().Write(tt.inputGps).Return(tt.returnErr).Times(tt.expectedCallTimes)
			eg.Go(func() error { return usecase.Execute(ctx, ch) })

			for i := 0; i < tt.expectedCallTimes; i++ {
				ch <- tt.inputGps
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
