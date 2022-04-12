package kitsunekkov2

import (
	"bytes"
	"gobot/pkg/animesubs"
	"gobot/pkg/stringutils"
	"gobot/pkg/webpagekeeper"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

type kitsunekkoScrapperV2 struct {
	logger          *zap.SugaredLogger
	timeUpdate time.Duration
	webkeeper webpagekeeper.WebPageKeeper
}

type pageEntry struct {
	text        string
	timeUpdated time.Time
	url         string
}

func (p pageEntry) Equal(other pageEntry) bool {
	return p.text == other.text &&
		p.timeUpdated.Equal(other.timeUpdated) &&
		p.url == other.url
}

var _ animesubs.AnimeSubsService = (*kitsunekkoScrapperV2)(nil)
var KitsunekkoTimeLayout = "Jan 02 2006 3:04:05 PM"
var kitsunekkoBaseUrl = "https://kitsunekko.net"
var kitsunekkoJapBaseUrl = "https://kitsunekko.net/dirlist.php?dir=subtitles%2Fjapanese%2F"

func parseKitsunekkoTime(timeString string) (time.Time, error) {
	parsedTime, err := time.Parse(KitsunekkoTimeLayout, timeString)
	if err != nil {
		return time.Time{}, err
	}
	return parsedTime, nil
}

func NewKitsunekkoScrapperV2(timeUpdate time.Duration, logger *zap.SugaredLogger) animesubs.AnimeSubsService {
	return &kitsunekkoScrapperV2 {
		logger: logger,
		timeUpdate: timeUpdate,
		webkeeper: webpagekeeper.NewWebPageKeeper(timeUpdate, logger),
	}
}

// TODO: should return error
func (serv *kitsunekkoScrapperV2) parseRowToEntry(s *goquery.Selection) pageEntry {
	if s.Children().Size() != 2 &&
		s.Children().Size() != 3 {
		return pageEntry{}
	}

	raw_url := s.Children().First()
	urlRaw, exist := raw_url.Find("a").Attr("href")
	if !exist {
		serv.logger.Errorf("Couldn't parse url of %v", raw_url.Text())
		return pageEntry{}
	}

	text := strings.TrimSpace(raw_url.Text())

	rawDate, exist := s.Children().Last().Attr("title")
	if !exist {
		serv.logger.Errorf("Couldn't parse raw date of %v", urlRaw)
		return pageEntry{}
	}

	parsedTime, err := parseKitsunekkoTime(rawDate)
	if err != nil {
		serv.logger.Errorf("Error in parsing %v string, error %v", rawDate, err)
		return pageEntry{}
	}


	return pageEntry{
		text: text,
		timeUpdated: parsedTime,
		url: urlRaw,
	}
}

func (serv *kitsunekkoScrapperV2) getAllEntries(body []byte) ([]pageEntry, error) {
	r := bytes.NewReader(body)

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	var entries []pageEntry
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		entry := serv.parseRowToEntry(s)
		if entry.url != "" {
			entries = append(entries, entry)
		}
	})

	return entries, nil
}

func (serv *kitsunekkoScrapperV2) filterPageEntriesByTitles(entries []pageEntry, titles []string) []pageEntry {
	var matchingEntries []pageEntry

	for _, titleRaw := range titles {
		title := stringutils.LowerAndTrimText(titleRaw)
		for _, entryRaw := range entries {
			entry := stringutils.LowerAndTrimText(entryRaw.text)
			if stringutils.GetLevenshteinDistancePercent(entry, title) < 80 {
				continue
			}
			matchingEntries = append(matchingEntries, entryRaw)
		}
	}
	return matchingEntries
}

func (serv *kitsunekkoScrapperV2) findLatestPageEntry(entries []pageEntry) pageEntry {
	if len(entries) == 0 {
		return pageEntry{}
	}

	actualentry := entries[0]
	latestTime := time.Unix(0, 0)

	for _, entry := range entries {
		if entry.timeUpdated.After(latestTime) {
			latestTime = entry.timeUpdated
			actualentry = entry
		}
	}

	return actualentry
}

func (serv *kitsunekkoScrapperV2) getRequiredAnimeUrl(titles []string) (string, error) {
	body, err := serv.webkeeper.GetUrlBody(kitsunekkoJapBaseUrl, true)
	if err != nil {
		return "", err
	}
	allEntries, err := serv.getAllEntries(body)
	if err != nil {
		return "", err
	}
	filteredEntries := serv.filterPageEntriesByTitles(allEntries, titles)
	lastEntry := serv.findLatestPageEntry(filteredEntries)
	return lastEntry.url, nil
}

func (serv *kitsunekkoScrapperV2) getLatestAnimeEntry(url string) (pageEntry, error) {
	body, err := serv.webkeeper.GetUrlBody(url, false)
	if err != nil {
		return pageEntry{}, err
	}

	allEntries, err := serv.getAllEntries(body)
	if err != nil {
		return pageEntry{}, err
	}
	lastEntry := serv.findLatestPageEntry(allEntries)
	return lastEntry, nil
}
 
func (serv *kitsunekkoScrapperV2) GetUrlLatestSubForAnime(titlesWithSynonyms []string) animesubs.SubsInfo {
	requiredUrl, err := serv.getRequiredAnimeUrl(titlesWithSynonyms)
	if err != nil {
		serv.logger.Errorf("Got error trying to get %v, error: %v", kitsunekkoJapBaseUrl, err)
		return animesubs.SubsInfo{}
	}

	if requiredUrl == "" {
		return animesubs.SubsInfo{}
	}

	entry, err := serv.getLatestAnimeEntry(kitsunekkoBaseUrl + requiredUrl)
	if err != nil {
		serv.logger.Errorf("Got error trying to get %v, error: %v", kitsunekkoBaseUrl + requiredUrl, err)
		return animesubs.SubsInfo{}
	}

	if entry.url == "" {
		serv.logger.Debugf("No entry was found for %v", titlesWithSynonyms[0])
		return animesubs.SubsInfo{}
	}
	
	if entry.url[0] != '/' {
		entry.url = "/" + entry.url
	}

	url, err := url.Parse(kitsunekkoBaseUrl + entry.url)
	if err != nil {
		serv.logger.Errorf("Broken url formed from %s and %s", kitsunekkoBaseUrl, entry.url)
		return animesubs.SubsInfo{}
	}

	return animesubs.SubsInfo{
		Title: entry.text,
		TimeUpdated: entry.timeUpdated,
		Url: url.String(),
	}
}