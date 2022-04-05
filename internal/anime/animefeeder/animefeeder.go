package animefeeder

import (
	"fmt"
	"gobot/internal/anime"
	"gobot/pkg/animeservice"
)

type AnimeFeeder interface {
	FeedInfo(info chan string)
}

type animeFeeder struct {
	animeService animeservice.AnimeService
	cachedList   *anime.AnimeList
}

var _ AnimeFeeder = (*animeFeeder)(nil)

func NewAnimeFeeder(animeService animeservice.AnimeService) AnimeFeeder {
	af := &animeFeeder{animeService: animeService, cachedList: anime.NewAnimeList()}
	af.cachedList.SetNewList(af.animeService.GetUserAnimeList())
	return af
}

func (af *animeFeeder) FeedInfo(info chan string) {
	defer close(info)
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
		info <- st
	}

	if missingInNew != nil {
		var st string
		st += "Entries were deleted\n"
		for _, v := range missingInNew {
			v.ListStatus = animeservice.NotInList
			st += v.VerboseOutput()
			st += "\n"
		}
		info <- st
	}

	af.cachedList.SetNewList(curList)

	fmt.Println("Feed info ended")
}
