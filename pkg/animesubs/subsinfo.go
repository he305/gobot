package animesubs

import "time"

type SubsInfo struct {
	Title       string
	TimeUpdated time.Time
	Url         string
}

func (s SubsInfo) IsEmpty() bool {
	return s.Title == "" || s.Url == ""
}

func (s SubsInfo) Equal(other SubsInfo) bool {
	return s.Title == other.Title &&
		s.TimeUpdated.Equal(other.TimeUpdated) &&
		s.Url == other.Url
}
