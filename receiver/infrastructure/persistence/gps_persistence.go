package persistence

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ktr03rtk/go-gps-logger/receiver/domain/model"
	"github.com/ktr03rtk/go-gps-logger/receiver/domain/repository"
	"github.com/pkg/errors"
)

type gpsPersistence struct {
	distDir string
}

func NewGpsPersistence(distDir string) (repository.GpsWriteRepository, error) {
	if err := os.MkdirAll(distDir, 0o755); err != nil {
		return nil, errors.Wrapf(err, "failed to create directors")
	}

	return &gpsPersistence{distDir}, nil
}

func (p *gpsPersistence) Write(g *model.Gps) error {
	filePath := createFilePath(p.distDir, g.Timestamp)

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		return errors.Wrapf(err, "failed to create file")
	}
	defer f.Close()

	bs, err := json.Marshal(g)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal json")
	}

	if _, err := fmt.Fprintln(f, string(bs)); err != nil {
		return errors.Wrapf(err, "failed to write file")
	}

	return nil
}

func createFilePath(distDir string, t time.Time) string {
	fileName := fmt.Sprintf("%s.dat", t.Format(("2006-01-02-15-04-05")))

	return filepath.Join(distDir, fileName)
}
