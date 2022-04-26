package animeurlservice_test

import (
	"gobot/pkg/animeurlservice"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsEmptyAnimeUrlInfo(t *testing.T) {
	testCases := []struct {
		Name     string
		Data     animeurlservice.AnimeUrlInfo
		Expected bool
	}{
		{
			"Empty title",
			animeurlservice.AnimeUrlInfo{
				Title:       "",
				Url:         "some url",
				TimeUpdated: time.Now(),
			},
			true,
		},
		{
			"Empty url",
			animeurlservice.AnimeUrlInfo{
				Title:       "some title",
				Url:         "",
				TimeUpdated: time.Now(),
			},
			true,
		},
		{
			"Both empty",
			animeurlservice.AnimeUrlInfo{
				Title:       "",
				Url:         "",
				TimeUpdated: time.Now(),
			},
			true,
		},
		{
			"Both present",
			animeurlservice.AnimeUrlInfo{
				Title:       "some title",
				Url:         "some url",
				TimeUpdated: time.Now(),
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
