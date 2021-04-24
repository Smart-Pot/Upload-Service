package service

import (
	"context"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"
)

type Service interface {
	UploadOne(ctx context.Context, file multipart.File) (string, error)
	UploadMany(ctx context.Context, file []multipart.File) ([]string, error)
}

func NewService(logg log.Logger) Service {
	wd, _ := os.Getwd()
	basedir := filepath.Join(wd, "upload")
	return &service{logg, basedir}
}

type service struct {
	logger  log.Logger
	basedir string
}

func (s service) UploadOne(ctx context.Context, file multipart.File) (filename string, err error) {
	defer func(beginTime time.Time) {
		level.Info(s.logger).Log(
			"function", "UploadOne",
			"result", filename,
			"err", err,
			"took", time.Since(beginTime))
	}(time.Now())
	filename = s.generateFileName()
	if err = s.saveFile(filename, file); err != nil {
		return "", err
	}

	return filename, err

}
func (s service) UploadMany(ctx context.Context, files []multipart.File) (filenames []string, err error) {
	defer func(beginTime time.Time) {
		level.Info(s.logger).Log(
			"function", "UploadMany",
			"result", filenames,
			"err", err,
			"took", time.Since(beginTime))
	}(time.Now())
	filenames = make([]string, len(files))
	for x := 0; x < len(files); x++ {
		filename := s.generateFileName()
		if err = s.saveFile(filename, files[x]); err != nil {
			return nil, err
		}
		filenames[x] = filename
	}

	return filenames, err
}

func (s service) generateFileName() string {
	str := uuid.NewString() + ".jpg"
	return str
}

func (s *service) saveFile(filename string, file multipart.File) error {
	defer file.Close()

	temp, err := os.Create(filepath.Join(s.basedir, filename))
	if err != nil {
		return err
	}
	defer temp.Close()

	fb, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	// write this byte array to our temporary file
	temp.Write(fb)
	return nil
}
