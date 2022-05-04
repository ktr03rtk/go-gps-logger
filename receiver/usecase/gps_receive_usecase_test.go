package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/golang/mock/gomock"
	"github.com/ktr03rtk/go-gps-logger/receiver/domain/model"
	"github.com/ktr03rtk/go-gps-logger/receiver/mock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func TestGpsReceiveUsecaseExecuteNormalCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repository := mock.NewMockGpsReceiveRepository(ctrl)
	mock := clock.NewMock()
	usecase := NewGpsReceiveUsecase(repository, mock)

	iteration := 2
	duration := 1000 * time.Millisecond
	date := time.Date(2022, 5, 3, 0, 9, 0, 0, time.Local)
	expectedOutput := &model.Gps{Timestamp: date, Mode: 3, Lat: 30.11, Lon: 130.11, Alt: 30.11, Speed: 30.11}
	var expectedErr error

	context, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(context)
	ch := make(chan *model.Gps)

	repository.EXPECT().Receive().Return(expectedOutput, expectedErr).Times(iteration)
	eg.Go(func() error { return usecase.Execute(ctx, duration, ch) })

	for i := 0; i < iteration; i++ {
		gosched()
		mock.Add(duration)

		output := <-ch
		assert.Exactly(t, expectedOutput, output)
	}

	cancel()

	if err := eg.Wait(); err != nil {
		t.Fatalf("error is not expected but received: %v", err)
	}
}

// Sleep momentarily so that other goroutines can process.
func gosched() { time.Sleep(1 * time.Millisecond) }
