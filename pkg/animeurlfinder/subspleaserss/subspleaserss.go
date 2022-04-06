package subspleaserss

import (
	"fmt"
	"gobot/pkg/animeurlfinder"
	"gobot/pkg/logging"
	"gobot/pkg/stringutils"
	"regexp"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"go.uber.org/zap"
)

var rss1080Url = "https://subsplease.org/rss/?t&r=1080"
var subspleaseRssPrefix = "[SubsPlease]"
var levenshteinPercentMin = 70

type subspleaserss struct {
	parser      *gofeed.Parser
	cachedFeed  *gofeed.Feed
	logger      *zap.SugaredLogger
	lastUpdated time.Time
}

var _ animeurlfinder.AnimeUrlFinder = (*subspleaserss)(nil)

func NewSubsPleaseRss() animeurlfinder.AnimeUrlFinder {
	return &subspleaserss{parser: gofeed.NewParser(), logger: logging.GetLogger(), cachedFeed: nil}
}

func (s *subspleaserss) updateFeed() {
	if s.cachedFeed == nil || time.Now().Sub(s.lastUpdated).Minutes() > 3 {
		feed, err := s.parser.ParseURL(rss1080Url)

		if err != nil {
			s.logger.Errorf("Could parse subs please rss url, url : %s, error: %s", rss1080Url, err.Error())
		}

		s.cachedFeed = feed
		s.lastUpdated = time.Now()
	}
}

func (s *subspleaserss) GetLatestUrlForTitle(title string) string {
	s.updateFeed()
	title = strings.TrimSpace(title)
	title = strings.ToLower(title)

	var rawRssTitles []string
	for _, rawTitle := range s.cachedFeed.Items {
		rawRssTitles = append(rawRssTitles, rawTitle.Title)
	}

	normalizedRssTitles := normalizeRssTitles(rawRssTitles)

	found := false
	var idx int
	for i, normalizedRssTitle := range normalizedRssTitles {
		if isRssMatchingTitle(normalizedRssTitle, title) {
			idx = i
			found = true
			break
		}
	}

	if found {
		return s.cachedFeed.Items[idx].Link
	}
	return ""
}

func isRssMatchingTitle(rss string, title string) bool {
	// First - simple matching
	isMatching := stringutils.AreSecondContainsFirst(title, rss)
	if isMatching {
		return true
	}

	// Second - levenshtein
	percent := stringutils.GetLevenshteinDistancePercent(rss, title)
	return percent >= levenshteinPercentMin
}

func normalizeRssTitles(titles []string) []string {
	var normalizeRssTitles []string

	re, err := regexp.Compile(`\[\w+\]`)

	// Bad, but will panic only if regexp module is upgraded
	if err != nil {
		panic(fmt.Sprintf("Could not compile regex in subsplease module, fatal error: %s", err.Error()))
	}

	for _, title := range titles {
		title = re.ReplaceAllString(title, "")
		title = strings.Replace(title, "(1080p)", "", -1)
		title = strings.Replace(title, ".mkv", "", -1)
		title = strings.ToLower(title)
		title = strings.Join(strings.Fields(title), " ")
		title = strings.TrimSpace(title)

		normalizeRssTitles = append(normalizeRssTitles, title)
	}

	return normalizeRssTitles
}
