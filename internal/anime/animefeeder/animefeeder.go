package animefeeder

import (
	"fmt"
	"gobot/internal/anime"
	"gobot/pkg/animeservice"
	"gobot/pkg/animesubs"
	"gobot/pkg/animeurlfinder"
)

type AnimeFeeder interface {
	UpdateList() (missingInCachedOutput []*animeservice.AnimeStruct, missingInNewOutput []*animeservice.AnimeStruct)
	FindLatestReleases() []LatestReleases
}

type LatestReleases struct {
	Anime    *animeservice.AnimeStruct
	AnimeUrl string
	SubsUrl  string
}

type animeFeeder struct {
	animeService   animeservice.AnimeService
	subServive     animesubs.AnimeSubsService
	animeUrlFinder animeurlfinder.AnimeUrlFinder
	cachedList     *anime.AnimeList
}

var _ AnimeFeeder = (*animeFeeder)(nil)

func NewAnimeFeeder(animeService animeservice.AnimeService, animesubs animesubs.AnimeSubsService, animeurlfinder animeurlfinder.AnimeUrlFinder) AnimeFeeder {
	af := &animeFeeder{animeService: animeService, cachedList: anime.NewAnimeList(), subServive: animesubs, animeUrlFinder: animeurlfinder}
	af.cachedList.SetNewList(af.animeService.GetUserAnimeList())
	return af
}

func (af *animeFeeder) UpdateList() (missingInCachedOutput []*animeservice.AnimeStruct, missingInNewOutput []*animeservice.AnimeStruct) {
	fmt.Println("Feed info called")
	curList := af.animeService.GetUserAnimeList()

	missingInCached, missingInNew := af.cachedList.FindMissingInBothLists(curList)

	af.cachedList.SetNewList(curList)

	fmt.Println("Feed info ended")

	missingInCachedOutput = missingInCached
	missingInNewOutput = missingInNew

	return
}

func (af *animeFeeder) FindLatestReleases() []LatestReleases {
	var releases []LatestReleases

	// Get filtered list
	filteredList := af.cachedList.FilterByListStatus(animeservice.PlannedToWatch, animeservice.Watching)

	for _, entry := range filteredList {
		// Check latest animeurl
		animeUrl := af.animeUrlFinder.GetLatestUrlForTitle(entry.Title)

		// Check latest subs
		animeSub := af.subServive.GetUrlLatestSubForAnime(entry.Title)

		if animeUrl != "" || animeSub != "" {
			releases = append(releases, LatestReleases{
				Anime:    entry,
				AnimeUrl: animeUrl,
				SubsUrl:  animeSub,
			})
		}
	}

	return releases
}
