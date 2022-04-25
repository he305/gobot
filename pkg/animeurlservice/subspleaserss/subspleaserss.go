package subspleaserss

import (
	"gobot/pkg/animeurlservice"
	"gobot/pkg/stringutils"
	"regexp"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"go.uber.org/zap"
)

const Rss1080Url = "https://subsplease.org/rss/?t&r=1080"

var rssSubsPleaseTimeLayout = "Mon, 02 Jan 2006 15:04:05 -0700"
var subspleaseRssPrefix = "[SubsPlease]"
var levenshteinPercentMin = 70

type subsPleaseRssEntry struct {
	Text        string
	RealText    string
	TimeUpdated time.Time
	Url         string
}

type subspleaserss struct {
	parser      *gofeed.Parser
	cachedFeed  *gofeed.Feed
	logger      *zap.SugaredLogger
	lastUpdated time.Time
	updateTimer time.Duration
	feedUrl     string
}

var _ animeurlservice.AnimeUrlService = (*subspleaserss)(nil)

func NewSubsPleaseRss(feedUrl string, updateTimer time.Duration, logger *zap.SugaredLogger) animeurlservice.AnimeUrlService {
	return &subspleaserss{
		parser:      gofeed.NewParser(),
		logger:      logger,
		cachedFeed:  nil,
		updateTimer: updateTimer,
		lastUpdated: time.Unix(0, 0),
		feedUrl:     feedUrl,
	}
}

func parseFeed(parser *gofeed.Parser, url string) (*gofeed.Feed, error) {
	return parser.ParseURL(url)
}

func (s *subspleaserss) updateFeed() error {
	if time.Now().Sub(s.lastUpdated) < s.updateTimer {
		return nil
	}

	feed, err := parseFeed(s.parser, s.feedUrl)

	if err != nil {
		return err
	}

	s.cachedFeed = feed
	s.lastUpdated = time.Now()

	return nil
}

func parseSubsPleaseTime(timeString string) (time.Time, error) {
	return time.Parse(rssSubsPleaseTimeLayout, timeString)
}

func parseSubsPleaseRssEntry(item *gofeed.Item) (subsPleaseRssEntry, error) {
	parsedTime, err := parseSubsPleaseTime(item.Published)
	if err != nil {
		return subsPleaseRssEntry{}, err
	}

	return subsPleaseRssEntry{
		Text:        item.Title,
		TimeUpdated: parsedTime,
		Url:         item.Link,
		RealText:    item.Title,
	}, nil
}

func (s *subspleaserss) getAllRssEntries() []subsPleaseRssEntry {
	var rawRss []subsPleaseRssEntry
	for _, item := range s.cachedFeed.Items {
		parsedEntry, err := parseSubsPleaseRssEntry(item)
		if err != nil {
			s.logger.Errorf("Error parsing entry %v, error: %v", item, err)
			continue
		}
		rawRss = append(rawRss, parsedEntry)
	}
	return rawRss
}

func (s *subspleaserss) getNormalizedRssEntries() []subsPleaseRssEntry {
	allEntries := s.getAllRssEntries()

	for i := range allEntries {
		allEntries[i].Text = normalizeRssTitle(allEntries[i].Text)
	}
	return allEntries
}

func (s *subspleaserss) filterEntriesByTitles(entries []subsPleaseRssEntry, titles ...string) []subsPleaseRssEntry {
	var filtered []subsPleaseRssEntry
	for _, entry := range entries {
		for _, title := range titles {
			title := stringutils.LowerAndTrimText(title)
			if isRssMatchingTitle(entry.Text, title) {
				s.logger.Infof("Found rss that matches title, title: %s, rss: %s", title, entry.Text)
				s.logger.Debugf("All titles: %v", titles)
				filtered = append(filtered, entry)
			}
		}
	}
	return filtered
}

func findLatestPageEntry(entries []subsPleaseRssEntry) subsPleaseRssEntry {
	if len(entries) == 0 {
		return subsPleaseRssEntry{}
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

func (s *subspleaserss) GetLatestUrlForTitle(titlesWithSynonyms ...string) animeurlservice.AnimeUrlInfo {
	if len(titlesWithSynonyms) == 0 || titlesWithSynonyms[0] == "" {
		return animeurlservice.AnimeUrlInfo{}
	}

	if err := s.updateFeed(); err != nil {
		s.logger.Errorf("Could parse subs please rss url, url : %s, error: %s", s.feedUrl, err.Error())
		return animeurlservice.AnimeUrlInfo{}
	}

	normalizedRssTitles := s.getNormalizedRssEntries()

	filteredEntries := s.filterEntriesByTitles(normalizedRssTitles, titlesWithSynonyms...)

	if len(filteredEntries) == 0 {
		return animeurlservice.AnimeUrlInfo{}
	}

	actualentry := findLatestPageEntry(filteredEntries)

	return animeurlservice.AnimeUrlInfo{
		Title:       actualentry.RealText,
		TimeUpdated: actualentry.TimeUpdated,
		Url:         actualentry.Url,
	}

}

func isRssMatchingTitle(rss string, title string) bool {
	// Normalize title
	title = stringutils.LowerAndTrimText(title)
	// First - simple matching
	isMatching := stringutils.AreSecondContainsFirst(title, rss)
	if isMatching {
		return true
	}

	// Second - levenshtein
	percent := stringutils.GetLevenshteinDistancePercent(rss, title)
	return percent >= levenshteinPercentMin
}

func normalizeRssTitle(title string) string {
	re := regexp.MustCompile(`\[\w+\]`)

	title = re.ReplaceAllString(title, "")
	title = strings.Replace(title, "(1080p)", "", -1)
	title = strings.Replace(title, ".mkv", "", -1)
	title = strings.ToLower(title)
	title = strings.Join(strings.Fields(title), " ")
	title = strings.TrimSpace(title)

	return title
}
