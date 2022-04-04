package animeservice

type AnimeService interface {
	GetAnimeByTitle(title string) *AnimeStruct
}
