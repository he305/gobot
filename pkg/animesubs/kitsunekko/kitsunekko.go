package kitsunekko

import (
	"fmt"
	"gobot/pkg/animesubs"
	"gobot/pkg/fileio"
	"gobot/pkg/logging"
	"gobot/pkg/stringutils"
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/storage"
	"go.uber.org/zap"
)

type kitsunekkoScrapper struct {
	logger          *zap.SugaredLogger
	lastTimeUpdated time.Time
	updateTimer     time.Duration
	cachedFilePath  string
	cachedVisitUrl  string
	fileIo          fileio.FileIO
	collector       *colly.Collector
}

type pageEntry struct {
	Text        string
	TimeUpdated time.Time
	Url         string
}

func (p pageEntry) Equal(other pageEntry) bool {
	return p.Text == other.Text &&
		p.TimeUpdated.Equal(other.TimeUpdated) &&
		p.Url == other.Url
}

var _ animesubs.AnimeSubsService = (*kitsunekkoScrapper)(nil)
var KitsunekkoTimeLayout = "Jan 02 2006 3:04:05 PM"
var kitsunekkoBaseUrl = "https://kitsunekko.net"
var kitsunekkoJapBaseUrl = "https://kitsunekko.net/dirlist.php?dir=subtitles%2Fjapanese%2F"

func getNewKitsunekkoCollyCollector() *colly.Collector {
	collector := colly.NewCollector()
	collector.Async = true
	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir(".")))
	collector.WithTransport(t)

	return collector
}

func NewKitsunekkoScrapper(fileIo fileio.FileIO, cachedFilePath string, updateTimer time.Duration) *kitsunekkoScrapper {
	return &kitsunekkoScrapper{logger: logging.GetLogger(),
		updateTimer:     updateTimer,
		cachedFilePath:  cachedFilePath,
		lastTimeUpdated: time.Unix(0, 0),
		cachedVisitUrl:  "file://" + cachedFilePath,
		fileIo:          fileIo,
		collector:       getNewKitsunekkoCollyCollector(),
	}
}

func (ws *kitsunekkoScrapper) processPageElement(e *colly.HTMLElement) (pageEntry, error) {

	// Just a random a[href], returning
	if e.DOM.Children().Size() == 0 {
		return pageEntry{}, nil
	}

	dateColumn := e.DOM.Parent().Siblings().Closest("td.tdright")

	timeString, ok := dateColumn.Attr("title")
	if !ok {
		return pageEntry{}, fmt.Errorf("Element with class td.tdright was not found, text: %s", e.Text)
	}

	parsedTime, err := parseKitsunekkoTime(timeString)
	if err != nil {
		return pageEntry{}, err
	}

	urlRaw := e.Request.AbsoluteURL(e.Attr("href"))

	return pageEntry{
		Text:        strings.TrimSpace(e.Text),
		TimeUpdated: parsedTime,
		Url:         urlRaw,
	}, nil
}

func filterPageEntriesByTitles(entries []pageEntry, titles []string) []pageEntry {
	var matchingEntries []pageEntry

	for _, titleRaw := range titles {
		title := stringutils.LowerAndTrimText(titleRaw)
		for _, entryRaw := range entries {
			entry := stringutils.LowerAndTrimText(entryRaw.Text)
			if stringutils.GetLevenshteinDistancePercent(entry, title) < 80 {
				continue
			}

			matchingEntries = append(matchingEntries, entryRaw)
		}
	}

	return matchingEntries
}

func findLatestPageEntry(entries []pageEntry) pageEntry {
	if len(entries) == 0 {
		return pageEntry{}
	}

	actualentry := entries[0]
	latestTime := time.Unix(0, 0)

	for _, entry := range entries {
		if entry.TimeUpdated.After(latestTime) {
			latestTime = entry.TimeUpdated
			actualentry = entry
		}
	}

	return actualentry
}

func (ws *kitsunekkoScrapper) getAllEntriesOnPage(path string) ([]pageEntry, error) {
	var allEntries []pageEntry

	collector := getNewKitsunekkoCollyCollector()

	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		parsedEntry, err := ws.processPageElement(e)
		if err != nil {
			ws.logger.Error(err)
			return
		}
		allEntries = append(allEntries, parsedEntry)
	})

	err := collector.Visit(path)
	collector.Wait()
	collector.SetStorage(&storage.InMemoryStorage{})
	return allEntries, err
}

func (ws *kitsunekkoScrapper) getLatestEntry(url string, titles ...string) pageEntry {
	allEntries, err := ws.getAllEntriesOnPage(url)
	if err != nil {
		ws.logger.Errorf("Error acquiring kitsunekko sub, url: %s, error: %s", url, err.Error())
	}

	if len(titles) != 0 {
		allEntries = filterPageEntriesByTitles(allEntries, titles)
	}

	actualEntry := findLatestPageEntry(allEntries)

	actualEntry.Url = strings.ReplaceAll(actualEntry.Url, "file://.", "")

	return actualEntry
}

func parseKitsunekkoTime(timeString string) (time.Time, error) {
	parsedTime, err := time.Parse(KitsunekkoTimeLayout, timeString)
	if err != nil {
		return time.Time{}, err
	}
	return parsedTime, nil
}

func (ws *kitsunekkoScrapper) GetUrlLatestSubForAnime(titlesWithSynonyms []string) animesubs.SubsInfo {
	mainPageUrl := ws.cachedVisitUrl

	err := ws.updateCache()
	if err != nil {
		ws.logger.Errorf("Error updating kitsunekko base site, error: %v", err)
		mainPageUrl = kitsunekkoJapBaseUrl
	}

	requiredUrl := ws.getLatestEntry(mainPageUrl, titlesWithSynonyms...)
	if requiredUrl.Url == "" {
		return animesubs.SubsInfo{}
	}

	time.Sleep(50 * time.Millisecond)
	actualEntry := ws.getLatestEntry(kitsunekkoBaseUrl + requiredUrl.Url)

	return animesubs.SubsInfo{
		Title:       actualEntry.Text,
		TimeUpdated: actualEntry.TimeUpdated,
		Url:         actualEntry.Url,
	}
}

func (ws *kitsunekkoScrapper) updateCache() error {
	if time.Now().Sub(ws.lastTimeUpdated) < ws.updateTimer {
		return nil
	}

	resp, err := http.Get(kitsunekkoJapBaseUrl)
	if err != nil {
		return err
	}

	err = ws.fileIo.SaveResponseToFile(resp, ws.cachedFilePath)
	if err != nil {
		return err
	}

	ws.lastTimeUpdated = time.Now()

	ws.logger.Infow("Kitsunekko index html was cached")
	return nil
}
