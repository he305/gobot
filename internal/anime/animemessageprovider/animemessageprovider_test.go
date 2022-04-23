package animemessageprovider

import (
	"fmt"
	"gobot/internal/anime/animefeeder"
	"gobot/pkg/animeservice"
	"gobot/pkg/animesubs"
	"gobot/pkg/animeurlservice"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockFeeder struct {
	BreakList bool
}

var _ animefeeder.AnimeFeeder = (*mockFeeder)(nil)

var (
	missingInCachedTest = []animeservice.AnimeStruct{
		{
			Title:        "shingeki no kyoujin",
			AiringStatus: 1,
			ListStatus:   0,
			ListRating:   0,
			ImageUrl:     "",
		},
		{
			Title:        "asobi asobase",
			AiringStatus: 1,
			ListStatus:   0,
			ListRating:   0,
			ImageUrl:     "",
		},
	}

	missingInNewTest = []animeservice.AnimeStruct{
		{
			Title:        "dom flandersa",
			AiringStatus: 1,
			ListStatus:   0,
			ListRating:   0,
			ImageUrl:     "",
		},
		{
			Title:        "spy x family",
			AiringStatus: 1,
			ListStatus:   0,
			ListRating:   0,
			ImageUrl:     "",
		},
	}
)

// Title: %s, airing status: %s, list status: %s, list rating: %d, image url:\n%s
func (m *mockFeeder) UpdateList() (missingInCachedOutput []animeservice.AnimeStruct, missingInNewOutput []animeservice.AnimeStruct, err error) {
	if m.BreakList {
		err = fmt.Errorf("BreakList")
		return
	}

	missingInCachedOutput = missingInCachedTest
	missingInNewOutput = missingInNewTest

	return
}

func (*mockFeeder) FindLatestReleases() []animefeeder.LatestReleases {
	return []animefeeder.LatestReleases{
		{
			AnimeUrl: animeurlservice.AnimeUrlInfo{
				Title:       "Shingeki no Kyoujin",
				Url:         "http://123",
				TimeUpdated: time.Now(),
			},
			SubsUrl: animesubs.SubsInfo{
				Title:       "Shingeki no Kyoujin",
				Url:         "http://12345",
				TimeUpdated: time.Now(),
			},
			Title: "Shingeki no Kyoujin",
		},
		{
			AnimeUrl: animeurlservice.AnimeUrlInfo{
				Title:       "Spy x family",
				Url:         "http://123",
				TimeUpdated: time.Now(),
			},
			SubsUrl: animesubs.SubsInfo{
				Title:       "Spy x family",
				Url:         "http://12345",
				TimeUpdated: time.Now(),
			},
			Title: "SPY X FAMILY",
		},
		{
			AnimeUrl: animeurlservice.AnimeUrlInfo{
				Title:       "asobi asobase",
				Url:         "http://123",
				TimeUpdated: time.Now(),
			},
			SubsUrl: animesubs.SubsInfo{
				Title:       "asobi asobase",
				Url:         "http://12345",
				TimeUpdated: time.Now(),
			},
			Title: "Asobi Asobase",
		},
	}
}

func newMockFeeder() animefeeder.AnimeFeeder {
	return &mockFeeder{}
}

func TestNewAnimeMessageProvider(t *testing.T) {
	assert := assert.New(t)

	expected := newMockFeeder()

	prov := NewAnimeMessageProvider(expected)

	assert.NotNil(prov)
}

func TestFormUpdateListMessage(t *testing.T) {
	assert := assert.New(t)

	feeder := newMockFeeder()
	prov := &animeMessageProvider{feeder: feeder}

	//Title: %s, airing status: %s, list status: %s, list rating: %d, image url:\n%s
	expected := "2 entries were deleted:\nTitle: dom flandersa, airing status: airing, list status: unknown, list rating: 0, image url:\n\nTitle: spy x family, airing status: airing, list status: unknown, list rating: 0, image url:\n\n2 entries were added:\nTitle: shingeki no kyoujin, airing status: airing, list status: unknown, list rating: 0, image url:\n\nTitle: asobi asobase, airing status: airing, list status: unknown, list rating: 0, image url:\n\n"

	actual, err := prov.formUpdatedListMessage()
	assert.NoError(err)
	assert.Equal(actual, expected)
}

func TestFormReleaseMessage(t *testing.T) {
	assert := assert.New(t)

	feeder := newMockFeeder()
	prov := &animeMessageProvider{feeder: feeder}
	expected := "New release for Shingeki no Kyoujin:\nNew series torrent:\nTitle: Shingeki no Kyoujin, Url: http://123\nNew series subs:\nTitle: Shingeki no Kyoujin, Url: http://12345\nNew release for SPY X FAMILY:\nNew series torrent:\nTitle: Spy x family, Url: http://123\nNew series subs:\nTitle: Spy x family, Url: http://12345\nNew release for Asobi Asobase:\nNew series torrent:\nTitle: asobi asobase, Url: http://123\nNew series subs:\nTitle: asobi asobase, Url: http://12345\n"

	actual := prov.formReleasesMessage()

	assert.Equal(expected, actual)
}

func TestGetMessage(t *testing.T) {
	assert := assert.New(t)

	feeder := newMockFeeder()
	prov := NewAnimeMessageProvider(feeder)

	expected := "2 entries were deleted:\nTitle: dom flandersa, airing status: airing, list status: unknown, list rating: 0, image url:\n\nTitle: spy x family, airing status: airing, list status: unknown, list rating: 0, image url:\n\n2 entries were added:\nTitle: shingeki no kyoujin, airing status: airing, list status: unknown, list rating: 0, image url:\n\nTitle: asobi asobase, airing status: airing, list status: unknown, list rating: 0, image url:\n\nNew release for Shingeki no Kyoujin:\nNew series torrent:\nTitle: Shingeki no Kyoujin, Url: http://123\nNew series subs:\nTitle: Shingeki no Kyoujin, Url: http://12345\nNew release for SPY X FAMILY:\nNew series torrent:\nTitle: Spy x family, Url: http://123\nNew series subs:\nTitle: Spy x family, Url: http://12345\nNew release for Asobi Asobase:\nNew series torrent:\nTitle: asobi asobase, Url: http://123\nNew series subs:\nTitle: asobi asobase, Url: http://12345\n"

	actual, err := prov.GetMessage()
	assert.NoError(err)
	assert.Equal(expected, actual)
}

func TestErrList(t *testing.T) {
	assert := assert.New(t)

	feeder := &mockFeeder{
		BreakList: true,
	}

	prov := NewAnimeMessageProvider(feeder)

	expected := ""
	actual, err := prov.GetMessage()
	assert.Error(err)
	assert.Equal(expected, actual)
}
