package kitsunekko

import (
	"fmt"
	"gobot/pkg/stringutils"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type KitsunekkoScrapper interface {
	GetUrlContainingString(url string, st string) string
}

type kitsunekkoScrapper struct {
	client *colly.Collector
}

var _ KitsunekkoScrapper = (*kitsunekkoScrapper)(nil)
var KitsunekkoTimeLayout = "Jan 02 2006 3:04:05 PM"

func NewKitsunekkoScrapper() KitsunekkoScrapper {
	return &kitsunekkoScrapper{client: colly.NewCollector()}
}

func (ws *kitsunekkoScrapper) GetUrlContainingString(url string, st string) string {
	type Entry struct {
		Title       string
		Url         string
		TimeUpdated time.Time
	}
	var founded []Entry

	ws.client.OnHTML("a[href]", func(e *colly.HTMLElement) {

		text := strings.ToLower(e.Text)
		if stringutils.GetLevenshteinDistancePercent(text, st) > 80 {
			timeSt, ok := e.DOM.Parent().Siblings().Attr("title")
			fmt.Println(timeSt)
			if !ok {
				return
			}

			parsedTime, err := time.Parse(KitsunekkoTimeLayout, timeSt)
			if err != nil {
				return
			}

			founded = append(founded, struct {
				Title       string
				Url         string
				TimeUpdated time.Time
			}{
				Title:       text,
				Url:         e.Request.AbsoluteURL(e.Attr("href")),
				TimeUpdated: parsedTime,
			})
		}
	})

	ws.client.Visit(url)

	actualEntry := founded[0]
	if len(founded) > 1 {
		latestTime := time.Unix(0, 0)
		for _, entry := range founded {
			if entry.TimeUpdated.After(latestTime) {
				latestTime = entry.TimeUpdated
				actualEntry = entry
			}
		}
	}

	fmt.Println(actualEntry)
	return actualEntry.Url
}
