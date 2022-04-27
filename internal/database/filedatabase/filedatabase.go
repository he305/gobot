package filedatabase

import (
	"encoding/json"
	"fmt"
	"gobot/internal/database"
	"gobot/pkg/fileio"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type fileDatabase struct {
	storageFolder string
	fileIo        fileio.FileIO
	logger        *zap.SugaredLogger
}

var _ database.Database = (*fileDatabase)(nil)

const DefaultExtension = "txt"

func NewFileDatabase(storageFolder string, logger *zap.SugaredLogger) *fileDatabase {
	storage := &fileDatabase{storageFolder: storageFolder, logger: logger, fileIo: fileio.NewDefaultFileIO()}
	return storage
}

func (db *fileDatabase) readFile(filepath string) (string, error) {
	content, err := db.fileIo.ReadFile(filepath)
	return string(content), err
}

func (db *fileDatabase) getFilesInDir(dirname string) ([]string, error) {
	return db.fileIo.GetFilesInDir(dirname)
}

func (db *fileDatabase) findFileNameInFolder(dirname string, name string) (string, error) {
	files, err := db.getFilesInDir(dirname)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		fileName := strings.ReplaceAll(file, "."+DefaultExtension, "")
		if fileName == name {
			return file, nil
		}
	}

	return "", nil
}

func (db *fileDatabase) GetEntryByName(collectionName string, key string, name string) (map[string]interface{}, error) {
	searchFolder := db.storageFolder + collectionName + "/"
	err := db.fileIo.CreateDirectory(searchFolder)
	if err != nil {
		return nil, err
	}

	filePath, err := db.findFileNameInFolder(searchFolder, name)
	if err != nil {
		return nil, err
	}

	if filePath == "" {
		return nil, nil
	}

	contents, err := db.fileIo.ReadFile(searchFolder + filePath)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}

	err = json.Unmarshal(contents, &result)
	return result, err
}

// AddEntry implements database.Database
func (db *fileDatabase) AddEntry(collectionName string, entry map[string]interface{}) error {
	if len(entry) == 0 {
		return nil
	}

	var fileName string

	if val, ok := entry["Title"]; ok {
		fileName = fmt.Sprintf("%v", val)
	} else {
		fileName = uuid.NewString()
	}

	folder := db.storageFolder + collectionName + "/"
	err := db.fileIo.CreateDirectory(folder)
	if err != nil {
		return err
	}

	fullPath := folder + fileName + "." + DefaultExtension

	jsonData, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	return db.fileIo.SaveToFile([]byte(jsonData), fullPath)
}

// AddEntries implements database.Database
func (db *fileDatabase) AddEntries(collectionName string, entries []map[string]interface{}) error {
	for _, entry := range entries {
		err := db.AddEntry(collectionName, entry)
		if err != nil {
			return err
		}
	}

	return nil
}
