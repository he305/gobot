package animesubs

type AnimeSubsService interface {
	GetUrlLatestSubForAnime(titlesWithSynonyms []string) SubsInfo
}
