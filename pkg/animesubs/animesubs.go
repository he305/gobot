package animesubs

type AnimeSubsService interface {
	GetUrlLatestSubForAnime(title string) string
}
