package usecase

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ktr03rtk/go-gps-logger/uploader/domain/model"
	"github.com/ktr03rtk/go-gps-logger/uploader/mock"
	"github.com/stretchr/testify/assert"
)

func TestFileReadUsecaseExecute(t *testing.T) {
	tests := []struct {
		name           string
		returnErr      error
		expectedOutput *model.Payload
		expectedErr    error
	}{
		{
			"normal case",
			nil,
			&model.Payload{Message: []byte("output-payload1\noutput-payload2"), FilePaths: []model.BaseFilePath{"file1", "file2"}},
			nil,
		},
		{
			"read error case",
			errors.New("error occurred"),
			&model.Payload{Message: []byte("output-payload1\noutput-payload2"), FilePaths: []model.BaseFilePath{"file1", "file2"}},
			errors.New("failed to execute file read usecase"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repository := mock.NewMockFileRepository(ctrl)
			usecase := NewFileReadUsecase(repository)

			repository.EXPECT().Read().Return(tt.expectedOutput, tt.returnErr).Times(1)

			output, err := usecase.Execute()
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
