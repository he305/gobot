package filereleasestorage

import (
	"gobot/internal/anime/animefeeder"
	"gobot/pkg/animesubs"
	"gobot/pkg/animeurlservice"
	"testing"
	"time"
)

func TestFormStringFromLatestReleases(t *testing.T) {
	expected := "Shingeki no kyoujin - 12" + defaultSeparator +
		"1649262976" + defaultSeparator +
		"sample.com/123" + defaultSeparator +
		"Shingeki no kyoujin - 12 [subs]" + defaultSeparator +
		"1649262976" + defaultSeparator +
		"https://anotherlink.com$%21"

	data := animefeeder.LatestReleases{
		AnimeUrl: animeurlservice.AnimeUrlInfo{
			Title:       "Shingeki no kyoujin - 12",
			TimeUpdated: time.Unix(1649262976, 0),
			Url:         "sample.com/123",
		},
		SubsUrl: animesubs.SubsInfo{
			Title:       "Shingeki no kyoujin - 12 [subs]",
			TimeUpdated: time.Unix(1649262976, 0),
			Url:         "https://anotherlink.com$%21",
		},
	}

	actual := formStringFromLatestReleases(data)

	if actual != expected {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}
