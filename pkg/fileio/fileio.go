package fileio

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type FileIO interface {
	SaveResponseToFile(data *http.Response, filePath string) error
	SaveToFile(data []byte, filePath string) error
	AppendToFile(data []byte, filePath string) error
	ReadFile(filePath string) ([]byte, error)
}

type fileIo struct{}

// ReadFile implements FileIO
func (*fileIo) ReadFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return content, nil
}

var _ FileIO = (*fileIo)(nil)

func NewDefaultFileIO() FileIO {
	return &fileIo{}
}

func writeToFile(file *os.File, data []byte) error {
	_, err := file.Write(data)
	return err
}

func (f *fileIo) SaveToFile(data []byte, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	return writeToFile(file, data)
}

func (f *fileIo) AppendToFile(data []byte, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	return writeToFile(file, data)
}

func (f *fileIo) SaveResponseToFile(data *http.Response, filePath string) error {
	defer data.Body.Close()
	byteData, err := ioutil.ReadAll(data.Body)
	if err != nil {
		return err
	}
	err = f.SaveToFile(byteData, filePath)
	return err
}
