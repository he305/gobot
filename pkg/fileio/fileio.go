package fileio

import (
	"io"
	"net/http"
	"os"
)

type FileIO interface {
	SaveResponseToFile(data *http.Response, filePath string) error
}

type fileIo struct{}

var _ FileIO = (*fileIo)(nil)

func NewDefaultFileIO() FileIO {
	return &fileIo{}
}

func (f *fileIo) SaveResponseToFile(data *http.Response, filePath string) error {
	defer data.Body.Close()
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, data.Body)
	if err != nil {
		return err
	}
	return nil
}
