package animeservice

import (
	"fmt"
	"time"
)

const (
	Airing uint8 = iota + 1
	CompletedAiring
	NotStarted
)

const (
	Unknown uint8 = iota
	NotInList
	Watching
	Completed
	PlannedToWatch
	Dropped
)

var mapAiringStatusToString = map[uint8]string{
	Airing:          "airing",
	CompletedAiring: "completed airing",
	NotStarted:      "not started",
}

var mapListStatusToString = map[uint8]string{
	Unknown:        "unknown",
	NotInList:      "not in list",
	Watching:       "watching",
	Completed:      "completed",
	PlannedToWatch: "planned to watch",
	Dropped:        "dropped",
}

type AnimeStruct struct {
	Id           int
	Title        string
	Synonyms     []string
	StartDate    time.Time
	EndTime      time.Time
	ListRating   float64
	AiringStatus uint8
	ListStatus   uint8
	ImageUrl     string
}

func NewAnimeStruct(id int, title string, synonyms []string, startDate time.Time, endTime time.Time, listRating float64, airingStatus uint8, listStatus uint8, imageUrl string) *AnimeStruct {
	return &AnimeStruct{Id: id, Title: title, Synonyms: synonyms, StartDate: startDate, ListRating: listRating, EndTime: endTime, AiringStatus: airingStatus, ListStatus: listStatus, ImageUrl: imageUrl}
}

func (a AnimeStruct) VerboseOutput() string {
	st := fmt.Sprintf("Title: %s, airing status: %s, list status: %s, list rating: %d, image url:\n%s", a.Title, mapAiringStatusToString[a.AiringStatus], mapListStatusToString[a.ListStatus], int(a.ListRating), a.ImageUrl)
	return st
}
