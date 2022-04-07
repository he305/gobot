package kitsunekko

import (
	"testing"
	"time"
)

// func TestGetUrlContainingName(t *testing.T) {
// 	expected := "https://kitsunekko.net/subtitles/japanese/Shingeki_No_Kyojin/Shingeki_No_Kyojin_097.srt"
// 	title := "shingeki no kyoujin"
// 	ws := kitsunekko.NewKitsunekkoScrapper()
// 	actual := ws.GetUrlLatestSubForAnime(title)
// 	if expected != actual {
// 		t.Errorf("expected %v, got %v", expected, actual)
// 	}
// }

func TestParseKitsunekkoTimeWithoutError(t *testing.T) {
	data := "Apr 05 2022 02:37:35 PM"
	expected := time.Date(2022, time.April, 5, 14, 37, 35, 0, time.UTC)

	got, err := parseKitsunekkoTime(data)

	if err != nil {
		t.Error("Error parsing kitsunekko time using default time layout")
	}

	if !got.Equal(expected) {
		t.Errorf("expected: %v, got: %v", expected, got)
	}
}

func TestParseKitsunekkoTimeWithError(t *testing.T) {
	data := "Apr 05 2022 02:37:35"
	expected := time.Time{}

	got, err := parseKitsunekkoTime(data)

	if err == nil {
		t.Error("Expected error while parsing kitsunekko time using default time layout")
	}

	if !got.Equal(expected) {
		t.Errorf("expected: %v, got: %v", expected, got)
	}
}

func TestFindLatestMainPageEntryEmptyArray(t *testing.T) {
	var data []pageEntry

	expected := pageEntry{}
	got := findLatestPageEntry(data)

	if got != expected {
		t.Errorf("expected: %v, got: %v", expected, got)
	}
}

func TestFindLatestMainPageEntryValid(t *testing.T) {
	data := []pageEntry{
		{
			Text:        "test1",
			TimeUpdated: time.Now(),
			Url:         "test",
		},
		{
			Text:        "test2",
			TimeUpdated: time.Now().Add(1 * time.Hour),
			Url:         "test",
		},
		{
			Text:        "test3",
			TimeUpdated: time.Now().Add(2 * time.Hour),
			Url:         "test",
		},
		{
			Text:        "test4",
			TimeUpdated: time.Now().Add(3 * time.Hour),
			Url:         "test",
		},
	}

	actual := data[len(data)-1]
	got := findLatestPageEntry(data)

	if !(actual.Text == got.Text &&
		actual.TimeUpdated.Equal(got.TimeUpdated) &&
		actual.Url == got.Url) {
		t.Errorf("expected: %v, got: %v", actual, got)
	}
}

func TestFilterMainPageEntriesByTitlesValid(t *testing.T) {
	entries := []pageEntry{
		{
			Text:        "THIS",
			TimeUpdated: time.Now(),
			Url:         "test",
		},
		{
			Text:        "tHiS",
			TimeUpdated: time.Now().Add(1 * time.Hour),
			Url:         "test",
		},
		{
			Text:        "test3",
			TimeUpdated: time.Now().Add(2 * time.Hour),
			Url:         "test",
		},
		{
			Text:        "test4",
			TimeUpdated: time.Now().Add(3 * time.Hour),
			Url:         "test",
		},
	}

	titles := []string{
		"irellevant",
		"THISA",
		"another text",
	}

	expected := []pageEntry{entries[0], entries[1]}
	got := filterPageEntriesByTitles(entries, titles)

	for _, expectedEn := range expected {
		found := false
		for _, gotEn := range got {
			if gotEn.Text == expectedEn.Text {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Entry %v not found in got list", expectedEn)
		}
	}
}

func TestConstructor(t *testing.T) {
	serv := NewKitsunekkoScrapper("", 1*time.Minute)
	if serv == nil {
		t.Errorf("Returned nil object")
	}
}

func TestCollyCollectorSettings(t *testing.T) {
	collector := configureKitsunekkoCollyCollector()
	if !collector.AllowURLRevisit ||
		len(collector.URLFilters) != 0 {
		t.Errorf("Wrong colly collector configuration")
	}
}
