package webpagekeeper

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type WebPageKeeper interface {
	GetUrlBody(url string, save bool) ([]byte, error)
}

type webPage struct {
	url string
	timeUpdated time.Time;
	body []byte
}

type webpagekeeper struct {
	timeToUpdate time.Duration;
	cachedWebPages []webPage
	logger *zap.SugaredLogger
}

var _ WebPageKeeper = (*webpagekeeper)(nil)

func NewWebPageKeeper(timeToUpdate time.Duration, logger *zap.SugaredLogger) WebPageKeeper {
	return &webpagekeeper{timeToUpdate: timeToUpdate, logger: logger}
}

func (wpk *webpagekeeper) getWebPage(url string) (webPage, error) {
	resp, err := http.Get(url)
	if err != nil {
		return webPage{}, err
	}

	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return webPage{}, fmt.Errorf("Couldn't get page %s, status code: %d", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return webPage{}, err
	}


	return webPage {
		url: url,
		body: body,
		timeUpdated: time.Now(),
	}, nil

}

func (wpk *webpagekeeper) GetUrlBody(url string, save bool) ([]byte, error) {
	now := time.Now()
	for i, cached := range wpk.cachedWebPages {
		if cached.url == url {
			if now.Sub(cached.timeUpdated) > wpk.timeToUpdate {
				newCached, err := wpk.getWebPage(cached.url)
				if err != nil {
					return nil, err
				}
				wpk.cachedWebPages[i] = newCached
				wpk.logger.Infof("%v was recached", url)
			}
			return cached.body, nil
		}
	}
	webpage, err := wpk.getWebPage(url)
	if err != nil {
		return nil, err
	}

	if save {
		wpk.cachedWebPages = append(wpk.cachedWebPages, webpage)
	}
	return webpage.body, nil
}