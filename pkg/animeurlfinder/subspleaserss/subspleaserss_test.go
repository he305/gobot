package subspleaserss

import (
	"testing"
	"time"
)

// func TestCorrectOuput(t *testing.T) {
// 	rss := NewSubsPleaseRss()

// 	expected := "https://nyaa.si/view/1510211/torrent"

// 	actual := rss.GetLatestUrlForTitle("Shingeki no Kyojin")
// 	fmt.Println(actual)

// 	if actual != expected {
// 		t.Errorf("expected: %v, got: %v", expected, actual)
// 	}
// }

func TestNormalizeRssTitles(t *testing.T) {
	data := []string{
		"[SubsPlease] Yami Shibai 10 - 13 (1080p) [A0A563BF].mkv",
		"[SubsPlease] Magia Record Final Season (09-12) (1080p) [Batch]",
		"[SubsPlease] Baraou no Souretsu - 12.5 (1080p) [19D386FA].mkv"}

	expected := []string{
		"yami shibai 10 - 13",
		"magia record final season (09-12)",
		"baraou no souretsu - 12.5",
	}

	actual := normalizeRssTitles(data)

	for i := 0; i < len(actual); i++ {
		if actual[i] != expected[i] {
			t.Errorf("expected %v, got %v", expected[i], actual[i])
		}
	}
}

func TestSubsPleaseTimeParse(t *testing.T) {
	data := "Tue, 05 Apr 2022 04:40:55 +0000"
	timeLayout := "Mon, 02 Jan 2006 15:04:05 -0700"

	expected := time.Date(2022, time.April, 5, 4, 40, 55, 0, time.UTC)
	actual, err := time.Parse(timeLayout, data)
	if err != nil {
		t.Errorf("expected %v, got %v", expected, err.Error())
	}

	if !actual.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}
