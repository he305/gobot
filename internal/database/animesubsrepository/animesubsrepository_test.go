package animesubsrepository

import (
	"gobot/internal/database"
	"gobot/pkg/animesubs"
	"gobot/pkg/animeurlservice"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type _mockDb struct {
}

// AddEntries implements database.Database
func (*_mockDb) AddEntries(collectionName string, entries []map[string]interface{}) error {
	return nil
}

// AddEntry implements database.Database
func (*_mockDb) AddEntry(collectionName string, entry map[string]interface{}) error {
	return nil
}

// GetEntryByName implements database.Database
func (*_mockDb) GetEntryByName(collectionName string, key string, name string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"title": name,
	}, nil
}

// AddAnimeUrls implements database.Database
func (*_mockDb) AddAnimeUrls(...animeurlservice.AnimeUrlInfo) error {
	return nil
}

// AddSubs implements database.Database
func (*_mockDb) AddSubs(...animesubs.SubsInfo) error {
	return nil
}

// GetAnimeSubByName implements database.Database
func (*_mockDb) GetAnimeSubByName(name string) (animesubs.SubsInfo, error) {
	return animesubs.SubsInfo{Title: name}, nil
}

// GetAnimeUrlByName implements database.Database
func (*_mockDb) GetAnimeUrlByName(name string) (animeurlservice.AnimeUrlInfo, error) {
	return animeurlservice.AnimeUrlInfo{}, nil
}

var _ database.Database = (*_mockDb)(nil)

func TestGetAnimeSubsByName(t *testing.T) {
	assert := assert.New(t)

	data := "testName"
	repo := NewAnimeSubsRepository(&_mockDb{}, "")
	expected := animesubs.SubsInfo{Title: data}

	actual, err := repo.GetAnimeSubsByName(data)
	assert.NoError(err)
	assert.Equal(expected.Title, actual.Title)
}

func TestAddAnimeSubsEmpty(t *testing.T) {
	assert := assert.New(t)

	repo := NewAnimeSubsRepository(&_mockDb{}, "")
	err := repo.AddAnimeSubs([]animesubs.SubsInfo{}...)
	assert.NoError(err)
}

var _okDataStructSubsInfo = []animesubs.SubsInfo{
	{
		Title:       "[Judas] Go-Toubun no Hanayome (Season 2).zip",
		Url:         "https://kitsunekko.net/subtitles/japanese/Gotoubun_No_Hanayome/%5BJudas%5D%20Go-Toubun%20no%20Hanayome%20%28Season%202%29.zip",
		TimeUpdated: time.Unix(1643792994, 0),
	},
	{
		Title:       "[Kunagisa] haibane renmei 03[h264.vorbis][cestfait].srt",
		TimeUpdated: time.Unix(1622283917, 0),
		Url:         "https://kitsunekko.net/subtitles/japanese/Haibane%20Renmei/%5BKunagisa%5D%20haibane_renmei_03%5Bh264.vorbis%5D%5Bcestfait%5D.srt",
	},
	{
		Title:       "JoJo no Kimyou na Bouken Stone Ocean [Netflix] 01.srt",
		TimeUpdated: time.Unix(1638405552, 0),
		Url:         "https://kitsunekko.net/subtitles/japanese/JoJo%20no%20Kimyou%20na%20Bouken:%20Stone%20Ocean/JoJo%20no%20Kimyou%20na%20Bouken%20Stone%20Ocean%20%5BNetflix%5D%2001.srt",
	},
}

var jsonSubs = []map[string]interface{}{
	{
		"Title": "[Judas] Go-Toubun no Hanayome (Season 2).zip",
		"Url":   "https://kitsunekko.net/subtitles/japanese/Gotoubun_No_Hanayome/%5BJudas%5D%20Go-Toubun%20no%20Hanayome%20%28Season%202%29.zip",
		"Time":  float64(1643792994),
	},
	{
		"Title": "[Kunagisa] haibane renmei 03[h264.vorbis][cestfait].srt",
		"Url":   "https://kitsunekko.net/subtitles/japanese/Haibane%20Renmei/%5BKunagisa%5D%20haibane_renmei_03%5Bh264.vorbis%5D%5Bcestfait%5D.srt",
		"Time":  float64(1622283917),
	},
	{
		"Title": "JoJo no Kimyou na Bouken Stone Ocean [Netflix] 01.srt",
		"Url":   "https://kitsunekko.net/subtitles/japanese/JoJo%20no%20Kimyou%20na%20Bouken:%20Stone%20Ocean/JoJo%20no%20Kimyou%20na%20Bouken%20Stone%20Ocean%20%5BNetflix%5D%2001.srt",
		"Time":  float64(1638405552),
	},
}

func TestAnimeSubsToJsonAnimeDataValid(t *testing.T) {
	assert := assert.New(t)

	type testStruct struct {
		name     string
		data     animesubs.SubsInfo
		expected map[string]interface{}
	}
	var testCases []testStruct

	assert.Equal(len(jsonSubs), len(_okDataStructSubsInfo))

	for i := range _okDataStructSubsInfo {
		testCases = append(testCases, testStruct{
			name:     _okDataStructSubsInfo[i].Title,
			data:     _okDataStructSubsInfo[i],
			expected: jsonSubs[i],
		})
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			animeData := formFileAnimeDataFromAnimeSubs(testCase.data)
			assert.False(animeData.IsEmpty())
			actual, err := formJsonFromAnimeData(animeData)
			assert.NoError(err)
			assert.Equal(testCase.expected, actual)
		})
	}
}

func TestAddAnimeSubsValid(t *testing.T) {
	assert := assert.New(t)

	repo := NewAnimeSubsRepository(&_mockDb{}, "")
	err := repo.AddAnimeSubs(_okDataStructSubsInfo...)
	assert.NoError(err)
}
