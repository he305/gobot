package animeurlservice

import "time"

type AnimeUrlInfo struct {
	Title       string
	TimeUpdated time.Time
	Url         string
}

func (a AnimeUrlInfo) Equal(other AnimeUrlInfo) bool {
	return a.Title == other.Title &&
		a.TimeUpdated.Equal(other.TimeUpdated) &&
		a.Url == other.Url
}

type AnimeUrlService interface {
	GetLatestUrlForTitle(titlesWithSynonyms ...string) AnimeUrlInfo
}
