package persistence

import (
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/ktr03rtk/go-gps-logger/uploader/domain/model"
	"github.com/ktr03rtk/go-gps-logger/uploader/domain/repository"
	"github.com/pkg/errors"
)

const fileExtention = ".raw"

type filePersistence struct {
	sourceDir string
}

func NewFilePersistence(sourceDir string) repository.FileRepository {
	return &filePersistence{sourceDir}
}

func (fp *filePersistence) Read() (*model.Payload, error) {
	payload := model.NewPayload()

	targetFiles, err := searchTargetFiles(fp.sourceDir)
	if err != nil {
		return nil, err
	}

	for _, file := range targetFiles {
		f, err := os.Open(file)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to open file")
		}
		defer f.Close()

		b, err := io.ReadAll(f)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read file")
		}

		payload.Add(b, model.BaseFilePath(filepath.Base(file)))
	}

	return payload, nil
}

func (fp *filePersistence) Delete(targetFiles []model.BaseFilePath) error {
	for _, baseFile := range targetFiles {
		if err := os.Remove(filepath.Join(fp.sourceDir, string(baseFile))); err != nil {
			return errors.Wrapf(err, "failed to remove file")
		}
	}

	return nil
}

func searchTargetFiles(sourceDir string) ([]string, error) {
	allFiles, err := filepath.Glob(sourceDir + "/*" + fileExtention)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find files")
	}

	sort.Strings(allFiles)
	targetFiles := allFiles

	if len(allFiles) >= model.MaxProcessFileNum {
		targetFiles = allFiles[0:model.MaxProcessFileNum]
	}

	return targetFiles, nil
}
