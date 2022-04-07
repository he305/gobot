package kitsunekko

import (
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

var _ animesubs.AnimeSubsService = (*kitsunekkoScrapper)(nil)
var KitsunekkoTimeLayout = "Jan 02 2006 3:04:05 PM"
var kitsunekkoBaseUrl = "https://kitsunekko.net"
var kitsunekkoJapBaseUrl = "https://kitsunekko.net/dirlist.php?dir=subtitles%2Fjapanese%2F"

func NewKitsunekkoScrapper(cachedFilePath string, updateTimer time.Duration) animesubs.AnimeSubsService {
	collector := colly.NewCollector()
	collector.AllowURLRevisit = true
	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir(".")))
	collector.WithTransport(t)

	return &kitsunekkoScrapper{logger: logging.GetLogger(), collector: collector, updateTimer: updateTimer, cachedFilePath: cachedFilePath, lastTimeUpdated: time.Unix(0, 0)}
}

func (ws *kitsunekkoScrapper) getRequiredAnimeUrl(titles []string) string {
	var founded []animesubs.SubsInfo

	ws.collector.OnHTML("a[href]", func(e *colly.HTMLElement) {

		text := strings.ToLower(e.Text)

		for _, title := range titles {
			if stringutils.GetLevenshteinDistancePercent(text, title) > 80 {
				timeSt, ok := e.DOM.Parent().Siblings().Attr("title")
				if !ok {
					return
				}

				parsedTime, err := time.Parse(KitsunekkoTimeLayout, timeSt)
				if err != nil {
					return
				}

				founded = append(founded, animesubs.SubsInfo{
					Title:       text,
					Url:         e.Attr("href"),
					TimeUpdated: parsedTime,
				})
			}
		}
	})

	if err := ws.collector.Visit("file://" + ws.cachedFilePath); err != nil {
		ws.logger.Errorf("Error acquiring kitsunekko sub, url: %s, error: %s", ws.cachedFilePath, err.Error())
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

	var en animesubs.SubsInfo
	latestTime := time.Unix(0, 0)
	ws.collector.OnHTML("td.tdright", func(e *colly.HTMLElement) {
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

			en = animesubs.SubsInfo{
				Title:       subTitle.Text(),
				TimeUpdated: parsedTime,
				Url:         e.Request.AbsoluteURL(localUrl),
			}
			latestTime = parsedTime
		}
	})

	// Let's sleep for some time before requesting second url
	time.Sleep(50 * time.Millisecond)
	if err := ws.collector.Visit(kitsunekkoBaseUrl + requiredUrl); err != nil {
		ws.logger.Errorf("Error acquiring kitsunekko sub, url: %s, error: %s", kitsunekkoBaseUrl+requiredUrl, err.Error())
		return animesubs.SubsInfo{}
	}

	return en
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
