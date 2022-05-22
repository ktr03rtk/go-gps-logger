//go:generate mockgen -source=file_repository.go -destination=../../mock/mock_file_repository.go -package=mock
package repository

import (
	"github.com/ktr03rtk/go-gps-logger/uploader/domain/model"
)

type FileRepository interface {
	Read() (*model.Payload, error)
	Delete([]model.BaseFilePath) error
}
