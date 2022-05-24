package usecase

import (
	"context"

	"github.com/ktr03rtk/go-gps-logger/uploader/domain/model"
	"github.com/ktr03rtk/go-gps-logger/uploader/domain/repository"
	"github.com/pkg/errors"
)

type PayloadUploadUsecase interface {
	Execute(context.Context, *model.Payload) ([]model.BaseFilePath, error)
}

type modelUploadUsecase struct {
	payloadRepository repository.PayloadUploadRepository
}

func NewPayloadUploadUsecase(pr repository.PayloadUploadRepository) PayloadUploadUsecase {
	return &modelUploadUsecase{
		payloadRepository: pr,
	}
}

func (uu *modelUploadUsecase) Execute(ctx context.Context, p *model.Payload) ([]model.BaseFilePath, error) {
	filePaths, err := uu.payloadRepository.Upload(ctx, p)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to execute payload upload usecase")
	}

	return filePaths, nil
}
