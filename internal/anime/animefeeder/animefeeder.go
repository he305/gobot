package animefeeder

import (
	"gobot/internal/anime"
	"gobot/internal/anime/animesubsrepository"
	"gobot/internal/anime/animeurlrepository"
	"gobot/pkg/animeservice"
	"gobot/pkg/animesubs"
	"gobot/pkg/animeurlservice"

	"go.uber.org/zap"
)

type AnimeFeeder interface {
	UpdateList() (missingInCachedOutput []animeservice.AnimeStruct, missingInNewOutput []animeservice.AnimeStruct, err error)
	FindLatestReleases() []LatestReleases
}

type LatestReleases struct {
	Title    string
	AnimeUrl animeurlservice.AnimeUrlInfo
	SubsUrl  animesubs.SubsInfo
}

func (l LatestReleases) Equal(other LatestReleases) bool {
	return l.AnimeUrl.Equal(other.AnimeUrl) &&
		l.SubsUrl.Equal(other.SubsUrl)
}

type animeFeeder struct {
	animeService        animeservice.AnimeService
	subServive          animesubs.AnimeSubsService
	animeUrlFinder      animeurlservice.AnimeUrlService
	animeUrlRepository  animeurlrepository.AnimeUrlRepository
	animeSubsRepository animesubsrepository.AnimeSubsRepository
	cachedList          *anime.AnimeList
	logger              *zap.SugaredLogger
	initialListError    bool
}

var _ AnimeFeeder = (*animeFeeder)(nil)

func NewAnimeFeeder(animeService animeservice.AnimeService,
	animesubs animesubs.AnimeSubsService,
	animeurlfinder animeurlservice.AnimeUrlService,
	animeUrlRepository animeurlrepository.AnimeUrlRepository,
	animeSubsRepository animesubsrepository.AnimeSubsRepository,
	logger *zap.SugaredLogger) AnimeFeeder {
	af := &animeFeeder{animeService: animeService,
		cachedList:          anime.NewAnimeList(),
		subServive:          animesubs,
		animeUrlFinder:      animeurlfinder,
		animeUrlRepository:  animeUrlRepository,
		animeSubsRepository: animeSubsRepository,
		logger:              logger}
	animeList, err := af.animeService.GetUserAnimeList()
	if err != nil {
		af.logger.Errorf("Error getting initial animelist")
		af.initialListError = true
	} else {
		af.initialListError = false
	}

	af.cachedList.SetNewList(animeList)
	return af
}

func (af *animeFeeder) UpdateList() (missingInCachedOutput []animeservice.AnimeStruct, missingInNewOutput []animeservice.AnimeStruct, err error) {
	curList, err := af.animeService.GetUserAnimeList()
	if err != nil {
		return nil, nil, err
	}

	if af.initialListError {
		af.cachedList.SetNewList(curList)
		af.initialListError = false
	}

	missingInCached, missingInNew := af.cachedList.FindMissingInBothLists(curList)

	af.cachedList.SetNewList(curList)

	missingInCachedOutput = missingInCached
	missingInNewOutput = missingInNew

	if len(missingInCachedOutput) != 0 || len(missingInNewOutput) != 0 {
		af.logger.Infof("Anime list was updated, %d entries added, %d entries deleted", len(missingInCachedOutput), len(missingInNewOutput))
	}
	return
}

func (af *animeFeeder) animeUrlExistInRepoOrNull(url animeurlservice.AnimeUrlInfo) bool {
	if url.Title == "" {
		return true
	}

	foundedAnimeUrl, err := af.animeUrlRepository.GetAnimeUrlByName(url.Title)
	if err != nil {
		af.logger.Errorf("Couldn't get info from anime repository, error: %v", err)
		return true
	}
	if foundedAnimeUrl.Url != "" {
		return true
	}

	return false
}

func (af *animeFeeder) animeSubsExistInRepoOrNull(subs animesubs.SubsInfo) bool {
	if subs.Title == "" {
		return true
	}

	foundedAnimeUrl, err := af.animeSubsRepository.GetAnimeSubsByName(subs.Title)
	if err != nil {
		af.logger.Errorf("Couldn't get info from subs repository, error: %v", err)
		return true
	}
	if foundedAnimeUrl.Url != "" {
		return true
	}

	return false
}

func (af *animeFeeder) FindLatestReleases() []LatestReleases {
	var releases []LatestReleases

	af.logger.Debug("Feeder finder started")
	// Get filtered list
	filteredList := af.cachedList.FilterByListStatus(animeservice.PlannedToWatch, animeservice.Watching)

	var newAnimeUrls []animeurlservice.AnimeUrlInfo
	var newSubs []animesubs.SubsInfo
	for _, entry := range filteredList {
		// Check latest animeurl
		animeUrlChan := make(chan animeurlservice.AnimeUrlInfo)
		go af.getLatestUrlForTitleChan(entry.FormAllNamesArray(), animeUrlChan)

		// Check latest subs
		animeSubChan := make(chan animesubs.SubsInfo)
		go af.getUrlLatestSubForAnimeChan(entry.FormAllNamesArray(), animeSubChan)

		animeUrl := <-animeUrlChan
		animeSub := <-animeSubChan

		isTrashAnimeUrl := af.animeUrlExistInRepoOrNull(animeUrl)

		if !isTrashAnimeUrl {
			newAnimeUrls = append(newAnimeUrls, animeUrl)
		} else {
			animeUrl = animeurlservice.AnimeUrlInfo{}
		}

		isTrashAnimeSubs := af.animeSubsExistInRepoOrNull(animeSub)

		if !isTrashAnimeSubs {
			newSubs = append(newSubs, animeSub)
		} else {
			animeSub = animesubs.SubsInfo{}
		}

		if !isTrashAnimeUrl || !isTrashAnimeSubs {
			releases = append(releases, LatestReleases{
				AnimeUrl: animeUrl,
				SubsUrl:  animeSub,
				Title:    entry.Title,
			})
		}
	}

	af.logger.Debug("Feeder finder ended")

	if err := af.animeUrlRepository.AddAnimeUrls(newAnimeUrls...); err != nil {
		af.logger.Errorf("Couldn't add %v anime urls to database, error: %v", len(newAnimeUrls), err)
	}

	if err := af.animeSubsRepository.AddAnimeSubs(newSubs...); err != nil {
		af.logger.Errorf("Couldn't add %v subs urls to database, error: %v", len(newSubs), err)
	}

	return releases
}

func (af *animeFeeder) getLatestUrlForTitleChan(titles []string, urlChan chan animeurlservice.AnimeUrlInfo) {
	data := af.animeUrlFinder.GetLatestUrlForTitle(titles...)
	urlChan <- data
	close(urlChan)
}

func (af *animeFeeder) getUrlLatestSubForAnimeChan(titles []string, subChan chan animesubs.SubsInfo) {
	data := af.subServive.GetUrlLatestSubForAnime(titles)
	subChan <- data
	close(subChan)
}
