package animesubs_test

import (
	"gobot/pkg/animesubs"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsEmptyAnimeUrlInfo(t *testing.T) {
	testCases := []struct {
		Name     string
		Data     animesubs.SubsInfo
		Expected bool
	}{
		{
			"Empty title",
			animesubs.SubsInfo{
				Title:       "",
				Url:         "some url",
				TimeUpdated: time.Now(),
			},
			true,
		},
		{
			"Empty url",
			animesubs.SubsInfo{
				Title:       "some title",
				Url:         "",
				TimeUpdated: time.Now(),
			},
			true,
		},
		{
			"Both empty",
			animesubs.SubsInfo{
				Title:       "",
				Url:         "",
				TimeUpdated: time.Now(),
			},
			true,
		},
		{
			"Both present",
			animesubs.SubsInfo{
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

func TestSubsInfoEqual(t *testing.T) {
	assert := assert.New(t)
	expected := true
	first := animesubs.SubsInfo{
		Title:       "test",
		TimeUpdated: time.Unix(0, 0),
		Url:         "http://example.com",
	}
	second := animesubs.SubsInfo{
		Title:       "test",
		TimeUpdated: time.Unix(0, 0),
		Url:         "http://example.com",
	}

	actual := first.Equal(second)
	assert.Equal(expected, actual)
}
