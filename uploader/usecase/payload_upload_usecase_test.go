package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ktr03rtk/go-gps-logger/uploader/domain/model"
	"github.com/ktr03rtk/go-gps-logger/uploader/mock"
	"github.com/stretchr/testify/assert"
)

func TestFileReadUsecaseUsecaseExecute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		input          *model.Payload
		returnErr      error
		expectedOutput []model.BaseFilePath
		expectedErr    error
	}{
		{
			"normal case",
			&model.Payload{Message: []byte("output-payload1\noutput-payload2"), FilePaths: []model.BaseFilePath{"file1", "file2"}},
			nil,
			[]model.BaseFilePath{"file1", "file2"},
			nil,
		},
		{
			"upload error case",
			&model.Payload{Message: []byte("output-payload1\noutput-payload2"), FilePaths: []model.BaseFilePath{"file1", "file2"}},
			errors.New("error occurred"),
			nil,
			errors.New("failed to execute payload upload usecase"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repository := mock.NewMockPayloadUploadRepository(ctrl)
			usecase := NewPayloadUploadUsecase(repository)

			ctx := context.Background()

			repository.EXPECT().Upload(ctx, tt.input).Return(tt.expectedOutput, tt.returnErr).Times(1)

			output, err := usecase.Execute(ctx, tt.input)
			if err != nil {
				if tt.expectedErr != nil {
					assert.Contains(t, err.Error(), tt.expectedErr.Error())
				} else {
					t.Fatalf("error is not expected but received: %v", err)
				}
			} else {
				assert.Exactly(t, tt.expectedErr, nil, "error is expected but received nil")
				assert.Exactly(t, tt.expectedOutput, output)
			}
		})
	}
}
