package mongodatabase

type AnimeEntry struct {
	Title string `bson:"title"`
	Url   string `bson:"url"`
	Time  int64  `bson:"time"`
}
