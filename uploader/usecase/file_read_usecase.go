package usecase

import (
	"github.com/ktr03rtk/go-gps-logger/uploader/domain/model"
	"github.com/ktr03rtk/go-gps-logger/uploader/domain/repository"
	"github.com/pkg/errors"
)

type FileReadUsecase interface {
	Execute() (*model.Payload, error)
}

type fileReadUsecase struct {
	fileRepository repository.FileRepository
}

func NewFileReadUsecase(fr repository.FileRepository) FileReadUsecase {
	return &fileReadUsecase{
		fileRepository: fr,
	}
}

func (fu *fileReadUsecase) Execute() (*model.Payload, error) {
	p, err := fu.fileRepository.Read()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to execute file read usecase")
	}

	return p, nil
}
