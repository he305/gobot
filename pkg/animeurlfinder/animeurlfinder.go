package animeurlfinder

import "time"

type AnimeUrlInfo struct {
	Title       string
	TimeUpdated time.Time
	Url         string
}

type AnimeUrlFinder interface {
	GetLatestUrlForTitle(titlesWithSynonyms []string) AnimeUrlInfo
}
