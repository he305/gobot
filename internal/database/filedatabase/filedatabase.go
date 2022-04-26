package filedatabase

import (
	"fmt"
	"gobot/internal/database"
	"gobot/pkg/animesubs"
	"gobot/pkg/animeurlservice"
	"gobot/pkg/fileio"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

type fileDatabase struct {
	animeUrlPathFile  string
	animeSubsPathFile string
	storagePath       string
	fileIo            fileio.FileIO
	logger            *zap.SugaredLogger
}

var _ database.Database = (*fileDatabase)(nil)
var defaultSeparator = "|"

func NewFileDatabase(animeUrlPathFile string, animesubsPathFile string, logger *zap.SugaredLogger) *fileDatabase {
	storage := &fileDatabase{animeUrlPathFile: animeUrlPathFile, animeSubsPathFile: animesubsPathFile, logger: logger, fileIo: fileio.NewDefaultFileIO()}
	return storage
}

func (db *fileDatabase) readFile(filepath string) (string, error) {
	content, err := db.fileIo.ReadFile(filepath)
	return string(content), err
}

func formAnimeData(data string) ([]FileAnimeData, error) {
	dataSplit := strings.Split(data, "\n")
	var fileAnimeDatas []FileAnimeData
	for _, line := range dataSplit {
		splitted := strings.Split(line, defaultSeparator)
		if len(splitted) < 3 {
			continue
		}
		rawTimeAnimeUrl, err := strconv.ParseInt(splitted[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Fatal error parsing %s line to unix time, error: %s", splitted[1], err.Error())
		}

		fileAnimeDatas = append(fileAnimeDatas, FileAnimeData{
			Title: splitted[0],
			Time:  rawTimeAnimeUrl,
			Url:   splitted[2],
		})
	}

	return fileAnimeDatas, nil
}

func formAnimeUrlInfos(data string) ([]animeurlservice.AnimeUrlInfo, error) {
	parsedData, err := formAnimeData(data)
	if err != nil {
		return nil, err
	}

	var savedAnimeUrls []animeurlservice.AnimeUrlInfo
	for _, entry := range parsedData {
		savedAnimeUrls = append(savedAnimeUrls, animeurlservice.AnimeUrlInfo{
			Title:       entry.Title,
			Url:         entry.Url,
			TimeUpdated: time.Unix(entry.Time, 0),
		})
	}

	return savedAnimeUrls, nil
}

func formAnimeSubs(data string) ([]animesubs.SubsInfo, error) {
	parsedData, err := formAnimeData(data)
	if err != nil {
		return nil, err
	}

	var savedSubs []animesubs.SubsInfo
	for _, entry := range parsedData {
		savedSubs = append(savedSubs, animesubs.SubsInfo{
			Title:       entry.Title,
			Url:         entry.Url,
			TimeUpdated: time.Unix(entry.Time, 0),
		})
	}

	return savedSubs, nil
}

func (db *fileDatabase) GetAnimeUrlByName(name string) (animeurlservice.AnimeUrlInfo, error) {
	data, err := db.readFile(db.animeUrlPathFile)
	if err != nil {
		return animeurlservice.AnimeUrlInfo{}, err
	}

	if data == "" {
		return animeurlservice.AnimeUrlInfo{}, nil
	}

	animeUrls, err := formAnimeUrlInfos(data)
	if err != nil {
		return animeurlservice.AnimeUrlInfo{}, err
	}

	for _, an := range animeUrls {
		if an.Title == name {
			return an, nil
		}
	}

	return animeurlservice.AnimeUrlInfo{}, nil
}

func (db *fileDatabase) GetAnimeSubByName(name string) (animesubs.SubsInfo, error) {
	data, err := db.readFile(db.animeSubsPathFile)
	if err != nil {
		return animesubs.SubsInfo{}, err
	}

	animeSubs, err := formAnimeSubs(data)
	if err != nil {
		return animesubs.SubsInfo{}, err
	}

	for _, an := range animeSubs {
		if an.Title == name {
			return an, nil
		}
	}
	return animesubs.SubsInfo{}, nil
}

func formFileAnimeDataFromAnimeUrl(entry animeurlservice.AnimeUrlInfo) FileAnimeData {
	return FileAnimeData{
		Title: entry.Title,
		Url:   entry.Url,
		Time:  entry.TimeUpdated.Unix(),
	}
}

func formFileAnimeDataFromAnimeSubs(entry animesubs.SubsInfo) FileAnimeData {
	return FileAnimeData{
		Title: entry.Title,
		Url:   entry.Url,
		Time:  entry.TimeUpdated.Unix(),
	}
}

func formStringDataToWrite(entry FileAnimeData) string {
	st := entry.Title + defaultSeparator +
		fmt.Sprintf("%d", entry.Time) + defaultSeparator +
		entry.Url
	return st
}

func (db *fileDatabase) writeToFile(data string, path string) error {
	return db.fileIo.AppendToFile([]byte(data), path)
}

func (db *fileDatabase) AddSubs(subs ...animesubs.SubsInfo) error {
	if len(subs) == 0 {
		return nil
	}

	var data string
	for _, entry := range subs {
		animeData := formFileAnimeDataFromAnimeSubs(entry)
		data += formStringDataToWrite(animeData) + "\n"
	}
	err := db.writeToFile(data, db.animeSubsPathFile)
	if err != nil {
		return err
	}

	db.logger.Infof("%d entries were saved into storage", len(subs))

	return nil
}

func (db *fileDatabase) AddAnimeUrls(urls ...animeurlservice.AnimeUrlInfo) error {
	if len(urls) == 0 {
		return nil
	}

	var data string
	for _, entry := range urls {
		animeData := formFileAnimeDataFromAnimeUrl(entry)
		data += formStringDataToWrite(animeData) + "\n"
	}
	err := db.writeToFile(data, db.animeUrlPathFile)
	if err != nil {
		return err
	}

	db.logger.Infof("%d entries were saved into storage", len(urls))

	return nil
}
