package entities_test

import (
	"gobot/internal/database/entities"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubsInfoEqual(t *testing.T) {
	assert := assert.New(t)
	expected := true
	first := entities.AnimeData{
		Title: "test",
		Time:  0,
		Url:   "http://example.com",
	}
	second := entities.AnimeData{
		Title: "test",
		Time:  0,
		Url:   "http://example.com",
	}

	actual := first.Equal(second)
	assert.Equal(expected, actual)
}

func TestIsEmptyAnimeUrlInfo(t *testing.T) {
	testCases := []struct {
		Name     string
		Data     entities.AnimeData
		Expected bool
	}{
		{
			"Empty title",
			entities.AnimeData{
				Title: "",
				Url:   "some url",
				Time:  0,
			},
			true,
		},
		{
			"Empty url",
			entities.AnimeData{
				Title: "some title",
				Url:   "",
				Time:  0,
			},
			true,
		},
		{
			"Both empty",
			entities.AnimeData{
				Title: "",
				Url:   "",
				Time:  0,
			},
			true,
		},
		{
			"Both present",
			entities.AnimeData{
				Title: "some title",
				Url:   "some url",
				Time:  0,
			},
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			assert.Equal(t, tc.Expected, tc.Data.IsEmpty())
		})
	}
}
