package filereleasestorage

import (
	"bufio"
	"fmt"
	"gobot/internal/anime/animefeeder"
	"gobot/internal/anime/releasestorage"
	"gobot/pkg/animesubs"
	"gobot/pkg/animeurlservice"
	"gobot/pkg/logging"
	"os"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

type fileReleaseStorage struct {
	cachedLatestRealeases []animefeeder.LatestReleases
	filePath              string
	logger                *zap.SugaredLogger
}

var _ releasestorage.ReleaseStorage = (*fileReleaseStorage)(nil)
var defaultSeparator = "|"

func NewFileReleaseStorage(path string) releasestorage.ReleaseStorage {
	storage := &fileReleaseStorage{filePath: path, logger: logging.GetLogger()}
	storage.readStorage()
	return storage
}

func (s *fileReleaseStorage) UpdateStorage(entries []animefeeder.LatestReleases) (newEntries []animefeeder.LatestReleases) {
	for _, entry := range entries {
		found := false
		for _, cachedEntry := range s.cachedLatestRealeases {
			if entry.Equal(cachedEntry) {
				found = true
			}
		}

		if !found {
			newEntries = append(newEntries, entry)
		}
	}

	s.cachedLatestRealeases = append(s.cachedLatestRealeases, newEntries...)
	s.saveToStorage(newEntries)

	return
}

func (s *fileReleaseStorage) readStorage() {
	f, err := os.Open(s.filePath)

	if err != nil {
		panic(fmt.Sprintf("Cannot open file %s, fatal error: %s", s.filePath, err.Error()))
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		splitted := strings.Split(line, defaultSeparator)

		rawTimeAnimeUrl, err := strconv.ParseInt(splitted[1], 10, 64)
		if err != nil {
			panic(fmt.Sprintf("Fatal error parsing %s line to unix time, error: %s", splitted[1], err.Error()))
		}
		parsedTimeAnimeUrl := time.Unix(rawTimeAnimeUrl, 0)

		rawTimeSubUrl, err := strconv.ParseInt(splitted[4], 10, 64)
		if err != nil {
			panic(fmt.Sprintf("Fatal error parsing %s line to unix time, error: %s", splitted[4], err.Error()))
		}
		parsedTimeSubUrl := time.Unix(rawTimeSubUrl, 0)

		s.cachedLatestRealeases = append(s.cachedLatestRealeases, animefeeder.LatestReleases{
			AnimeUrl: animeurlservice.AnimeUrlInfo{
				Title:       splitted[0],
				TimeUpdated: parsedTimeAnimeUrl,
				Url:         splitted[2],
			},
			SubsUrl: animesubs.SubsInfo{
				Title:       splitted[3],
				TimeUpdated: parsedTimeSubUrl,
				Url:         splitted[5],
			},
		})
	}
	s.logger.Infow("File storage is read")
}

func (s *fileReleaseStorage) saveToStorage(newEntries []animefeeder.LatestReleases) {
	if len(newEntries) == 0 {
		return
	}

	f, err := os.OpenFile(s.filePath, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(fmt.Sprintf("Couldn't open file %s, error %s", s.filePath, err.Error()))
	}
	defer f.Close()

	for _, entry := range newEntries {
		st := formStringFromLatestReleases(entry)
		st += "\n"

		if _, err := f.WriteString(st); err != nil {
			panic(fmt.Sprintf("Couldn't write %s to file %s, error: %s", st, s.filePath, err.Error()))
		}
	}

	s.logger.Infof("%d entries were saved into storage", len(newEntries))
}

func formStringFromLatestReleases(entry animefeeder.LatestReleases) string {
	st := entry.AnimeUrl.Title + defaultSeparator +
		fmt.Sprintf("%d", entry.AnimeUrl.TimeUpdated.Unix()) + defaultSeparator +
		entry.AnimeUrl.Url + defaultSeparator +
		entry.SubsUrl.Title + defaultSeparator +
		fmt.Sprintf("%d", entry.SubsUrl.TimeUpdated.Unix()) + defaultSeparator +
		entry.SubsUrl.Url
	return st
}
