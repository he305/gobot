package animesubs

import "time"

type SubsInfo struct {
	Title       string
	TimeUpdated time.Time
	Url         string
}
type AnimeSubsService interface {
	GetUrlLatestSubForAnime(titlesWithSynonyms []string) SubsInfo
}
