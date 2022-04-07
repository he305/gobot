package animefeeder

import (
	"gobot/internal/anime"
	"gobot/pkg/animeservice"
	"gobot/pkg/animesubs"
	"gobot/pkg/animeurlfinder"
	"gobot/pkg/logging"

	"go.uber.org/zap"
)

type AnimeFeeder interface {
	UpdateList() (missingInCachedOutput []*animeservice.AnimeStruct, missingInNewOutput []*animeservice.AnimeStruct)
	FindLatestReleases() []LatestReleases
}

type LatestReleases struct {
	Anime    *animeservice.AnimeStruct
	AnimeUrl animeurlfinder.AnimeUrlInfo
	SubsUrl  animesubs.SubsInfo
}

func (l LatestReleases) Equal(other LatestReleases) bool {
	return l.AnimeUrl.Equal(other.AnimeUrl) &&
		l.SubsUrl.Equal(other.SubsUrl)
}

type animeFeeder struct {
	animeService   animeservice.AnimeService
	subServive     animesubs.AnimeSubsService
	animeUrlFinder animeurlfinder.AnimeUrlFinder
	cachedList     *anime.AnimeList
	logger         *zap.SugaredLogger
}

var _ AnimeFeeder = (*animeFeeder)(nil)

func NewAnimeFeeder(animeService animeservice.AnimeService, animesubs animesubs.AnimeSubsService, animeurlfinder animeurlfinder.AnimeUrlFinder) AnimeFeeder {
	af := &animeFeeder{animeService: animeService, cachedList: anime.NewAnimeList(), subServive: animesubs, animeUrlFinder: animeurlfinder, logger: logging.GetLogger()}
	af.cachedList.SetNewList(af.animeService.GetUserAnimeList())
	return af
}

func (af *animeFeeder) UpdateList() (missingInCachedOutput []*animeservice.AnimeStruct, missingInNewOutput []*animeservice.AnimeStruct) {
	curList := af.animeService.GetUserAnimeList()

	missingInCached, missingInNew := af.cachedList.FindMissingInBothLists(curList)

	af.cachedList.SetNewList(curList)

	missingInCachedOutput = missingInCached
	missingInNewOutput = missingInNew

	if len(missingInCachedOutput) != 0 || len(missingInNewOutput) != 0 {
		af.logger.Infof("Anime list was updated, %d entries added, %d entries deleted", len(missingInCachedOutput), len(missingInNewOutput))
	}
	return
}

func (af *animeFeeder) FindLatestReleases() []LatestReleases {
	var releases []LatestReleases

	// Get filtered list
	filteredList := af.cachedList.FilterByListStatus(animeservice.PlannedToWatch, animeservice.Watching)

	for _, entry := range filteredList {
		// Check latest animeurl
		animeUrlChan := make(chan animeurlfinder.AnimeUrlInfo)
		go af.getLatestUrlForTitleChan(entry.FormAllNamesArray(), animeUrlChan)

		// Check latest subs
		animeSubChan := make(chan animesubs.SubsInfo)
		go af.getUrlLatestSubForAnimeChan(entry.FormAllNamesArray(), animeSubChan)

		animeUrl := <-animeUrlChan
		animeSub := <-animeSubChan

		if animeUrl.Url != "" || animeSub.Url != "" {
			releases = append(releases, LatestReleases{
				Anime:    entry,
				AnimeUrl: animeUrl,
				SubsUrl:  animeSub,
			})
		}
	}

	return releases
}

func (af *animeFeeder) getLatestUrlForTitleChan(titles []string, urlChan chan animeurlfinder.AnimeUrlInfo) {
	data := af.animeUrlFinder.GetLatestUrlForTitle(titles)
	urlChan <- data
	close(urlChan)
}

func (af *animeFeeder) getUrlLatestSubForAnimeChan(titles []string, subChan chan animesubs.SubsInfo) {
	data := af.subServive.GetUrlLatestSubForAnime(titles)
	subChan <- data
	close(subChan)
}
