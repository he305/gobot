package animeurlservice

type AnimeUrlService interface {
	GetLatestUrlForTitle(titlesWithSynonyms ...string) AnimeUrlInfo
}
