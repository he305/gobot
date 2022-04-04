package animeservice

import "time"

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

type AnimeStruct struct {
	Title        string
	Synonyms     []string
	StartDate    time.Time
	EndTime      time.Time
	ListRating   float64
	AiringStatus uint8
	ListStatus   uint8
}

func NewAnimeStruct(title string, synonyms []string, startDate time.Time, endTime time.Time, listRating float64, airingStatus uint8, listStatus uint8) *AnimeStruct {
	return &AnimeStruct{Title: title, Synonyms: synonyms, StartDate: startDate, EndTime: endTime, ListStatus: listStatus}
}
