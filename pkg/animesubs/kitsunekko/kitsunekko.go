package kitsunekko

import (
	"fmt"
	"gobot/pkg/animesubs"
	"gobot/pkg/logging"
	"gobot/pkg/stringutils"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"go.uber.org/zap"
)

type kitsunekkoScrapper struct {
	logger          *zap.SugaredLogger
	collector       *colly.Collector
	lastTimeUpdated time.Time
	updateTimer     time.Duration
	cachedFilePath  string
}

type pageEntry struct {
	Text        string
	TimeUpdated time.Time
	Url         string
}

var _ animesubs.AnimeSubsService = (*kitsunekkoScrapper)(nil)
var KitsunekkoTimeLayout = "Jan 02 2006 3:04:05 PM"
var kitsunekkoBaseUrl = "https://kitsunekko.net"
var kitsunekkoJapBaseUrl = "https://kitsunekko.net/dirlist.php?dir=subtitles%2Fjapanese%2F"

func configureKitsunekkoCollyCollector() *colly.Collector {
	collector := colly.NewCollector()
	collector.AllowURLRevisit = true
	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir(".")))
	collector.WithTransport(t)

	return collector
}

func NewKitsunekkoScrapper(cachedFilePath string, updateTimer time.Duration) animesubs.AnimeSubsService {
	return &kitsunekkoScrapper{logger: logging.GetLogger(), collector: configureKitsunekkoCollyCollector(), updateTimer: updateTimer, cachedFilePath: cachedFilePath, lastTimeUpdated: time.Unix(0, 0)}
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
		Text:        e.Text,
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

	ws.collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		parsedEntry, err := ws.processPageElement(e)
		if err != nil {
			ws.logger.Error(err)
			return
		}
		allEntries = append(allEntries, parsedEntry)
	})

	err := ws.collector.Visit(path)
	return allEntries, err
}

func (ws *kitsunekkoScrapper) getRequiredAnimeUrl(titles []string) string {
	allEntries, err := ws.getAllEntriesOnPage("file://" + ws.cachedFilePath)
	if err != nil {
		ws.logger.Errorf("Error acquiring kitsunekko sub, url: %s, error: %s", ws.cachedFilePath, err.Error())
	}

	matchingEntries := filterPageEntriesByTitles(allEntries, titles)

	actualEntry := findLatestPageEntry(matchingEntries)

	fixedUrl := strings.ReplaceAll(actualEntry.Url, "file://.", "")
	return fixedUrl
}

func parseKitsunekkoTime(timeString string) (time.Time, error) {
	parsedTime, err := time.Parse(KitsunekkoTimeLayout, timeString)
	if err != nil {
		return time.Time{}, err
	}
	return parsedTime, nil
}

func (ws *kitsunekkoScrapper) GetUrlLatestSubForAnime(titlesWithSynonyms []string) animesubs.SubsInfo {
	err := ws.updateCache()
	if err != nil {
		ws.logger.Errorf("Error updating kitsunekko base site, error: %v", err)
		return animesubs.SubsInfo{}
	}

	requiredUrl := ws.getRequiredAnimeUrl(titlesWithSynonyms)
	if requiredUrl == "" {
		return animesubs.SubsInfo{}
	}

	allEntries, err := ws.getAllEntriesOnPage(kitsunekkoBaseUrl + requiredUrl)
	if err != nil {
		ws.logger.Errorf("Error acquiring kitsunekko sub, url: %s, error: %s", ws.cachedFilePath, err.Error())
		return animesubs.SubsInfo{}
	}

	actualEntry := findLatestPageEntry(allEntries)

	return animesubs.SubsInfo{
		Title:       actualEntry.Text,
		TimeUpdated: actualEntry.TimeUpdated,
		Url:         actualEntry.Url,
	}
}

func (ws *kitsunekkoScrapper) updateCache() error {
	if time.Now().Sub(ws.lastTimeUpdated) >= ws.updateTimer {
		resp, err := http.Get(kitsunekkoJapBaseUrl)
		if err != nil {
			return err
		}

		defer resp.Body.Close()
		f, err := os.OpenFile(ws.cachedFilePath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(f, resp.Body)
		if err != nil {
			return err
		}
		ws.lastTimeUpdated = time.Now()

		ws.logger.Infow("Kitsunekko index html was cached")
	}

	return nil
}
