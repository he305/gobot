package animeservice_test

import (
	"gobot/pkg/animeservice"
	"testing"
	"time"
)

func TestVerboseAnimeStruct(t *testing.T) {
	anime := animeservice.NewAnimeStruct(
		0,
		"Sample",
		nil,
		time.Now(),
		time.Now(),
		6,
		animeservice.Airing,
		animeservice.Dropped,
	)

	expectedSt := "Title: Sample, airing status: airing, list status: dropped, list rating: 6"

	actualSt := anime.VerboseOutput()

	if actualSt != expectedSt {
		t.Errorf("got %s, wanted %s", actualSt, expectedSt)
	}
}
