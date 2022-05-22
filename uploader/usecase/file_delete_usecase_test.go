package usecase

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ktr03rtk/go-gps-logger/uploader/domain/model"
	"github.com/ktr03rtk/go-gps-logger/uploader/mock"
	"github.com/stretchr/testify/assert"
)

func TestFileDeleteUsecaseExecute(t *testing.T) {
	tests := []struct {
		name        string
		input       []model.BaseFilePath
		returnErr   error
		expectedErr error
	}{
		{
			"normal case",
			[]model.BaseFilePath{"file1", "file2"},
			nil,
			nil,
		},
		{
			"delete error case",
			[]model.BaseFilePath{"file1", "file2"},
			errors.New("error occurred"),
			errors.New("failed to execute file delete usecase"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repository := mock.NewMockFileRepository(ctrl)
			usecase := NewFileDeleteUsecase(repository)

			repository.EXPECT().Delete(tt.input).Return(tt.returnErr).Times(1)

			err := usecase.Execute(tt.input)
			if err != nil {
				if tt.expectedErr != nil {
					assert.Contains(t, err.Error(), tt.expectedErr.Error())
				} else {
					t.Fatalf("error is not expected but received: %v", err)
				}
			} else {
				assert.Exactly(t, tt.expectedErr, nil, "error is expected but received nil")
				assert.Nil(t, err)
			}
		})
	}
}
