package filedatabase

type FileAnimeData struct {
	Title string
	Time  int64
	Url   string
}

func (f FileAnimeData) Equal(other FileAnimeData) bool {
	return f.Title == other.Title &&
		f.Time == other.Time &&
		f.Url == other.Url
}
