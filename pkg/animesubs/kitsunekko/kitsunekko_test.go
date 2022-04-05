package kitsunekko_test

import (
	"gobot/pkg/animesubs/kitsunekko"
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

func TestFormatTimeCorrect(t *testing.T) {
	acutalTime := time.Date(2022, time.April, 5, 14, 37, 35, 0, time.UTC)
	time, err := time.Parse(kitsunekko.KitsunekkoTimeLayout, "Apr 05 2022 02:37:35 PM")
	if err != nil {
		t.Error("Error parsing kitsunekko time using default time layout")
	}

	if time != acutalTime {
		t.Errorf("expected: %v, got: %v", acutalTime, time)
	}

}
