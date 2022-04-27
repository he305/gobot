package entities

var defaultSeparator = "|"

type AnimeData struct {
	Title string
	Time  int64
	Url   string
}

func (f AnimeData) Equal(other AnimeData) bool {
	return f.Title == other.Title &&
		f.Time == other.Time &&
		f.Url == other.Url
}

func (f AnimeData) IsEmpty() bool {
	return f.Title == "" || f.Url == ""
}
