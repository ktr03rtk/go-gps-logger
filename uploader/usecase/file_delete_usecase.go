package usecase

import (
	"github.com/ktr03rtk/go-gps-logger/uploader/domain/model"
	"github.com/ktr03rtk/go-gps-logger/uploader/domain/repository"
	"github.com/pkg/errors"
)

type FileDeleteUsecase interface {
	Execute([]model.BaseFilePath) error
}

type fileDeleteUsecase struct {
	fileRepository repository.FileRepository
}

func NewFileDeleteUsecase(fr repository.FileRepository) FileDeleteUsecase {
	return &fileDeleteUsecase{
		fileRepository: fr,
	}
}

func (fu *fileDeleteUsecase) Execute(filePaths []model.BaseFilePath) error {
	if err := fu.fileRepository.Delete(filePaths); err != nil {
		return errors.Wrapf(err, "failed to execute file delete usecase")
	}

	return nil
}
