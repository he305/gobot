package animeservice

import "time"

type airingStatusType int8

const (
	Airing airingStatusType = iota + 1
	CompletedAiring
	NotStarted
)

type listStatusType int8

const (
	Unknown listStatusType = iota
	NotInList
	Watching
	Completed
	PlannedToWatch
	Dropped
)

type AnimeStruct struct {
	title        string
	synonyms     []string
	startDate    time.Time
	endTime      time.Time
	globalRating float64
	listRating   float64
	airingStatus airingStatusType
	listStatus   listStatusType
}

func NewAnimeStruct(title string, synonyms []string, startDate time.Time, endTime time.Time, globalRating float64, listRating float64, airingStatus airingStatusType, listStatus listStatusType) *AnimeStruct {
	return &AnimeStruct{title: title, synonyms: synonyms, startDate: startDate, endTime: endTime, globalRating: globalRating, listStatus: listStatus}
}
