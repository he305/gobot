package kitsunekko

import (
	"gobot/pkg/animesubs"
	"gobot/pkg/logging"
	"gobot/pkg/stringutils"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"go.uber.org/zap"
)

type kitsunekkoScrapper struct {
	logger *zap.SugaredLogger
}

type entry struct {
	Title       string
	Url         string
	TimeUpdated time.Time
}

var _ animesubs.AnimeSubsService = (*kitsunekkoScrapper)(nil)
var KitsunekkoTimeLayout = "Jan 02 2006 3:04:05 PM"
var kitsunekkoJapBaseUrl = "https://kitsunekko.net/dirlist.php?dir=subtitles%2Fjapanese%2F"

func NewKitsunekkoScrapper() animesubs.AnimeSubsService {
	return &kitsunekkoScrapper{logger: logging.GetLogger()}
}

func (ws *kitsunekkoScrapper) getRequiredAnimeUrl(title string) string {
	var founded []entry

	collector := colly.NewCollector()

	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {

		text := strings.ToLower(e.Text)
		if stringutils.GetLevenshteinDistancePercent(text, title) > 80 {
			timeSt, ok := e.DOM.Parent().Siblings().Attr("title")
			if !ok {
				return
			}

			parsedTime, err := time.Parse(KitsunekkoTimeLayout, timeSt)
			if err != nil {
				return
			}

			founded = append(founded, entry{
				Title:       text,
				Url:         e.Request.AbsoluteURL(e.Attr("href")),
				TimeUpdated: parsedTime,
			})
		}
	})

	if err := collector.Visit(kitsunekkoJapBaseUrl); err != nil {
		ws.logger.Errorf("Error acquiring kitsunekko sub, url: %s, error: %s", kitsunekkoJapBaseUrl, err.Error())
		return ""
	}

	if len(founded) == 0 {
		return ""
	}

	actualentry := founded[0]
	if len(founded) > 1 {
		latestTime := time.Unix(0, 0)
		for _, entry := range founded {
			if entry.TimeUpdated.After(latestTime) {
				latestTime = entry.TimeUpdated
				actualentry = entry
			}
		}
	}

	return actualentry.Url
}

func (ws *kitsunekkoScrapper) GetUrlLatestSubForAnime(title string) string {
	requiredUrl := ws.getRequiredAnimeUrl(title)
	if requiredUrl == "" {
		return ""
	}

	collector := colly.NewCollector()

	var en entry
	latestTime := time.Unix(0, 0)
	collector.OnHTML("td.tdright", func(e *colly.HTMLElement) {
		timeSt := e.Attr("title")

		parsedTime, err := time.Parse(KitsunekkoTimeLayout, timeSt)
		if err != nil {
			return
		}

		if parsedTime.After(latestTime) {

			subTitle := e.DOM.Siblings().Find("a[href]")
			localUrl, exist := subTitle.Attr("href")
			if !exist {
				return
			}

			en = entry{
				Title:       subTitle.Text(),
				TimeUpdated: parsedTime,
				Url:         e.Request.AbsoluteURL(localUrl),
			}
			latestTime = parsedTime
		}
	})

	// Let's sleep for some time before requesting second url
	time.Sleep(100 * time.Millisecond)
	if err := collector.Visit(requiredUrl); err != nil {
		ws.logger.Errorf("Error acquiring kitsunekko sub, url: %s, error: %s", requiredUrl, err.Error())
		return ""
	}

	return en.Url
}
