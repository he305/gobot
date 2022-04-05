package kitsunekko_test

import (
	"gobot/internal/anime/kitsunekko"
	"testing"
	"time"
)

// func TestGetUrlContainingName(t *testing.T) {
// 	ws := kitsunekko.NewKitsunekkoScrapper()
// 	ws.GetUrlContainingString("https://kitsunekko.net/dirlist.php?dir=subtitles%2Fjapanese%2F", "shingeki no kyoujin")
// 	// t.Log(len(actual))
// 	// for a := range test {
// 	// 	t.Log(a)
// 	// }
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
