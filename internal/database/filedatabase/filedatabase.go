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

// func (s *fileDatabase) UpdateStorage(entries []animefeeder.LatestReleases) (newEntries []animefeeder.LatestReleases) {
// 	for _, entry := range entries {
// 		found := false
// 		for _, cachedEntry := range s.cachedAnimeUrls {
// 			if entry.Equal(cachedEntry) {
// 				found = true
// 			}
// 		}

// 		if !found {
// 			newEntries = append(newEntries, entry)
// 		}
// 	}

// 	s.cachedAnimeUrls = append(s.cachedAnimeUrls, newEntries...)
// 	s.saveToStorage(newEntries)

// 	return
// }

// func (s *fileDatabase) readStorage() {
// 	f, err := os.Open(s.storagePath)

// 	if err != nil {
// 		panic(fmt.Sprintf("Cannot open file %s, fatal error: %s", s.storagePath, err.Error()))
// 	}

// 	defer f.Close()

// 	scanner := bufio.NewScanner(f)
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		splitted := strings.Split(line, defaultSeparator)

// 		rawTimeAnimeUrl, err := strconv.ParseInt(splitted[1], 10, 64)
// 		if err != nil {
// 			panic(fmt.Sprintf("Fatal error parsing %s line to unix time, error: %s", splitted[1], err.Error()))
// 		}
// 		parsedTimeAnimeUrl := time.Unix(rawTimeAnimeUrl, 0)

// 		rawTimeSubUrl, err := strconv.ParseInt(splitted[4], 10, 64)
// 		if err != nil {
// 			panic(fmt.Sprintf("Fatal error parsing %s line to unix time, error: %s", splitted[4], err.Error()))
// 		}
// 		parsedTimeSubUrl := time.Unix(rawTimeSubUrl, 0)

// 		s.cachedAnimeUrls = append(s.cachedAnimeUrls, animefeeder.LatestReleases{
// 			AnimeUrl: animeurlservice.AnimeUrlInfo{
// 				Title:       splitted[0],
// 				TimeUpdated: parsedTimeAnimeUrl,
// 				Url:         splitted[2],
// 			},
// 			SubsUrl: animesubs.SubsInfo{
// 				Title:       splitted[3],
// 				TimeUpdated: parsedTimeSubUrl,
// 				Url:         splitted[5],
// 			},
// 		})
// 	}
// 	s.logger.Infow("File storage is read")
// }

// func (s *fileDatabase) saveToStorage(newEntries []animefeeder.LatestReleases) {
// 	if len(newEntries) == 0 {
// 		return
// 	}

// 	f, err := os.OpenFile(s.storagePath, os.O_APPEND|os.O_WRONLY, 0666)
// 	if err != nil {
// 		panic(fmt.Sprintf("Couldn't open file %s, error %s", s.storagePath, err.Error()))
// 	}
// 	defer f.Close()

// 	for _, entry := range newEntries {
// 		st := formStringFromLatestReleases(entry)
// 		st += "\n"

// 		if _, err := f.WriteString(st); err != nil {
// 			panic(fmt.Sprintf("Couldn't write %s to file %s, error: %s", st, s.storagePath, err.Error()))
// 		}
// 	}

// 	s.logger.Infof("%d entries were saved into storage", len(newEntries))
// }

// func formStringFromLatestReleases(entry animefeeder.LatestReleases) string {
// 	st := entry.AnimeUrl.Title + defaultSeparator +
// 		fmt.Sprintf("%d", entry.AnimeUrl.TimeUpdated.Unix()) + defaultSeparator +
// 		entry.AnimeUrl.Url + defaultSeparator +
// 		entry.SubsUrl.Title + defaultSeparator +
// 		fmt.Sprintf("%d", entry.SubsUrl.TimeUpdated.Unix()) + defaultSeparator +
// 		entry.SubsUrl.Url
// 	return st
// }
