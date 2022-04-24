package filedatabase

import (
	"fmt"
	"gobot/internal/database"
	"gobot/pkg/animesubs"
	"gobot/pkg/animeurlservice"
	"gobot/pkg/logging"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

type fileDatabase struct {
	animeUrlPathFile  string
	animeSubsPathFile string
	storagePath       string
	logger            *zap.SugaredLogger
}

var _ database.Database = (*fileDatabase)(nil)
var defaultSeparator = "|"

func NewFileDatabase(animeUrlPathFile string, animesubsPathFile string) database.Database {
	storage := &fileDatabase{animeUrlPathFile: animeUrlPathFile, animeSubsPathFile: animesubsPathFile, logger: logging.GetLogger()}
	return storage
}

func readFile(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func formAnimeUrlInfos(data string) []animeurlservice.AnimeUrlInfo {
	dataSplit := strings.Split(data, "\n")
	var savedAnimeUrls []animeurlservice.AnimeUrlInfo
	for _, line := range dataSplit {
		splitted := strings.Split(line, defaultSeparator)
		if len(splitted) < 3 {
			continue
		}
		rawTimeAnimeUrl, err := strconv.ParseInt(splitted[1], 10, 64)
		if err != nil {
			panic(fmt.Sprintf("Fatal error parsing %s line to unix time, error: %s", splitted[1], err.Error()))
		}
		parsedTimeAnimeUrl := time.Unix(rawTimeAnimeUrl, 0)

		savedAnimeUrls = append(savedAnimeUrls, animeurlservice.AnimeUrlInfo{
			Title:       splitted[0],
			TimeUpdated: parsedTimeAnimeUrl,
			Url:         splitted[2],
		})
	}

	return savedAnimeUrls
}

func formAnimeSubs(data string) []animesubs.SubsInfo {
	dataSplit := strings.Split(data, "\n")
	var savedSubs []animesubs.SubsInfo
	for _, line := range dataSplit {
		splitted := strings.Split(line, defaultSeparator)
		if len(splitted) < 3 {
			continue
		}

		rawTimeAnimeUrl, err := strconv.ParseInt(splitted[1], 10, 64)
		if err != nil {
			panic(fmt.Sprintf("Fatal error parsing %s line to unix time, error: %s", splitted[1], err.Error()))
		}
		parsedTimeAnimeUrl := time.Unix(rawTimeAnimeUrl, 0)

		savedSubs = append(savedSubs, animesubs.SubsInfo{
			Title:       splitted[0],
			TimeUpdated: parsedTimeAnimeUrl,
			Url:         splitted[2],
		})
	}

	return savedSubs
}

func (db *fileDatabase) GetAnimeUrlByName(name string) (animeurlservice.AnimeUrlInfo, error) {
	data, err := readFile(db.animeUrlPathFile)
	if err != nil {
		return animeurlservice.AnimeUrlInfo{}, err
	}

	if data == "" {
		return animeurlservice.AnimeUrlInfo{}, nil
	}

	animeUrls := formAnimeUrlInfos(data)

	for _, an := range animeUrls {
		if an.Title == name {
			return an, nil
		}
	}

	return animeurlservice.AnimeUrlInfo{}, nil
}

func (db *fileDatabase) AddAnimeUrls(urls ...animeurlservice.AnimeUrlInfo) error {
	if len(urls) == 0 {
		return nil
	}

	f, err := os.OpenFile(db.animeUrlPathFile, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, entry := range urls {
		st := formAnimeUrlForWrite(entry)
		st += "\n"

		if _, err := f.WriteString(st); err != nil {
			return fmt.Errorf("Couldn't write %s to file %s, error: %s", st, db.animeUrlPathFile, err.Error())
		}
	}

	db.logger.Infof("%d entries were saved into storage", len(urls))

	return nil
}

func formAnimeUrlForWrite(entry animeurlservice.AnimeUrlInfo) string {
	st := entry.Title + defaultSeparator +
		fmt.Sprintf("%d", entry.TimeUpdated.Unix()) + defaultSeparator +
		entry.Url
	return st
}

func formSubsForWrite(entry animesubs.SubsInfo) string {
	st := entry.Title + defaultSeparator +
		fmt.Sprintf("%d", entry.TimeUpdated.Unix()) + defaultSeparator +
		entry.Url
	return st
}

func (db *fileDatabase) GetAnimeSubByName(name string) (animesubs.SubsInfo, error) {
	data, err := readFile(db.animeSubsPathFile)
	if err != nil {
		return animesubs.SubsInfo{}, err
	}

	animeSubs := formAnimeSubs(data)

	for _, an := range animeSubs {
		if an.Title == name {
			return an, nil
		}
	}
	return animesubs.SubsInfo{}, nil
}

func (db *fileDatabase) AddSubs(subs ...animesubs.SubsInfo) error {
	if len(subs) == 0 {
		return nil
	}

	f, err := os.OpenFile(db.animeSubsPathFile, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, entry := range subs {
		st := formSubsForWrite(entry)
		st += "\n"

		if _, err := f.WriteString(st); err != nil {
			return fmt.Errorf("Couldn't write %s to file %s, error: %s", st, db.animeSubsPathFile, err.Error())
		}
	}

	db.logger.Infof("%d entries were saved into storage", len(subs))

	return nil
}
