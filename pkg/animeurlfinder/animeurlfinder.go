package animeurlfinder

type AnimeUrlFinder interface {
	GetLatestUrlForTitle(title string) string
}
