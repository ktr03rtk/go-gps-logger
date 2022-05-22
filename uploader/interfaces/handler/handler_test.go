package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/golang/mock/gomock"
	"github.com/ktr03rtk/go-gps-logger/uploader/domain/model"
	"github.com/ktr03rtk/go-gps-logger/uploader/mock"
	"github.com/ktr03rtk/go-gps-logger/uploader/usecase"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func TestLogUploadHandler(t *testing.T) {
	t.Parallel()
	duration := 1000 * time.Millisecond

	tests := []struct {
		name            string
		readOutput      *model.Payload
		readErr         error
		uploadOutput    []model.BaseFilePath
		uploadErr       error
		deleteErr       error
		expectedErr     error
		readCallTimes   int
		uploadCallTimes int
		deleteCallTimes int
	}{
		{
			"normal case",
			&model.Payload{Message: []byte("output-payload1\noutput-payload2"), FilePaths: []model.BaseFilePath{"file1", "file2"}},
			nil,
			[]model.BaseFilePath{"file1", "file2"},
			nil,
			nil,
			nil,
			1,
			1,
			1,
		},
		{
			"read error case",
			&model.Payload{Message: []byte("output-payload1\noutput-payload2"), FilePaths: []model.BaseFilePath{"file1", "file2"}},
			errors.New("error occurred"),
			nil,
			nil,
			nil,
			errors.New("error occurred"),
			1,
			0,
			0,
		},
		{
			"upload error case",
			&model.Payload{Message: []byte("output-payload1\noutput-payload2"), FilePaths: []model.BaseFilePath{"file1", "file2"}},
			nil,
			[]model.BaseFilePath{"file1", "file2"},
			errors.New("error occurred"),
			nil,
			errors.New("error occurred"),
			1,
			1,
			0,
		},
		{
			"delete error case",
			&model.Payload{Message: []byte("output-payload1\noutput-payload2"), FilePaths: []model.BaseFilePath{"file1", "file2"}},
			nil,
			[]model.BaseFilePath{"file1", "file2"},
			nil,
			errors.New("error occurred"),
			errors.New("error occurred"),
			1,
			1,
			1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fr := mock.NewMockFileRepository(ctrl)
			pr := mock.NewMockPayloadUploadRepository(ctrl)
			ru := usecase.NewFileReadUsecase(fr)
			uu := usecase.NewPayloadUploadUsecase(pr)
			du := usecase.NewFileDeleteUsecase(fr)
			mock := clock.NewMock()
			h := NewLogUploadHandler(ru, uu, du, mock)

			fr.EXPECT().Read().Return(tt.readOutput, tt.readErr).Times(tt.readCallTimes)
			pr.EXPECT().Upload(tt.readOutput).Return(tt.uploadOutput, tt.uploadErr).Times(tt.uploadCallTimes)
			fr.EXPECT().Delete(tt.uploadOutput).Return(tt.deleteErr).Times(tt.deleteCallTimes)

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
