//go:generate mockgen -source=payload_upload_repository.go -destination=../../mock/mock_payload_upload_repository.go -package=mock
package repository

import (
	"github.com/ktr03rtk/go-gps-logger/uploader/domain/model"
)

type PayloadUploadRepository interface {
	Upload(*model.Payload) ([]model.BaseFilePath, error)
}
