package animefeeder

import (
	"fmt"
	"gobot/internal/anime"
	"gobot/pkg/animeservice"
	"gobot/pkg/animesubs"
)

type AnimeFeeder interface {
	UpdateList() (missingInCachedOutput []*animeservice.AnimeStruct, missingInNewOutput []*animeservice.AnimeStruct)
}

type animeFeeder struct {
	animeService animeservice.AnimeService
	subServive   animesubs.AnimeSubsService
	cachedList   *anime.AnimeList
}

var _ AnimeFeeder = (*animeFeeder)(nil)

func NewAnimeFeeder(animeService animeservice.AnimeService, animesubs animesubs.AnimeSubsService) AnimeFeeder {
	af := &animeFeeder{animeService: animeService, cachedList: anime.NewAnimeList(), subServive: animesubs}
	af.cachedList.SetNewList(af.animeService.GetUserAnimeList())
	return af
}

func (af *animeFeeder) UpdateList() (missingInCachedOutput []*animeservice.AnimeStruct, missingInNewOutput []*animeservice.AnimeStruct) {
	fmt.Println("Feed info called")
	curList := af.animeService.GetUserAnimeList()

	missingInCached, missingInNew := af.cachedList.FindMissingInBothLists(curList)

	fmt.Println(len(missingInCached), len(missingInNew))

	if missingInCached != nil {
		var st string
		st += "New entries in list\n"
		for _, v := range missingInCached {
			st += v.VerboseOutput()
			st += "\n"
		}
	}

	if missingInNew != nil {
		var st string
		st += "Entries were deleted\n"
		for _, v := range missingInNew {
			v.ListStatus = animeservice.NotInList
			st += v.VerboseOutput()
			st += "\n"
		}
	}

	af.cachedList.SetNewList(curList)

	fmt.Println("Feed info ended")

	missingInCachedOutput = missingInCached
	missingInNewOutput = missingInNew

	return
}
