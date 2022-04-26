package filedatabase

import (
	"gobot/pkg/animesubs"
	"gobot/pkg/animeurlservice"
	"gobot/pkg/fileio"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var (
	_okDataSubsFilePath = "oksubs"
	_okDataSubs         = "[Judas] Go-Toubun no Hanayome (Season 2).zip|1643792994|https://kitsunekko.net/subtitles/japanese/Gotoubun_No_Hanayome/%5BJudas%5D%20Go-Toubun%20no%20Hanayome%20%28Season%202%29.zip\n[Kunagisa] haibane renmei 03[h264.vorbis][cestfait].srt|1622283917|https://kitsunekko.net/subtitles/japanese/Haibane%20Renmei/%5BKunagisa%5D%20haibane_renmei_03%5Bh264.vorbis%5D%5Bcestfait%5D.srt\nJoJo no Kimyou na Bouken Stone Ocean [Netflix] 01.srt|1638405552|https://kitsunekko.net/subtitles/japanese/JoJo%20no%20Kimyou%20na%20Bouken:%20Stone%20Ocean/JoJo%20no%20Kimyou%20na%20Bouken%20Stone%20Ocean%20%5BNetflix%5D%2001.srt\n"
	_okDataStructSubs   = []FileAnimeData{
		{
			Title: "[Judas] Go-Toubun no Hanayome (Season 2).zip",
			Url:   "https://kitsunekko.net/subtitles/japanese/Gotoubun_No_Hanayome/%5BJudas%5D%20Go-Toubun%20no%20Hanayome%20%28Season%202%29.zip",
			Time:  1643792994,
		},
		{
			Title: "[Kunagisa] haibane renmei 03[h264.vorbis][cestfait].srt",
			Time:  1622283917,
			Url:   "https://kitsunekko.net/subtitles/japanese/Haibane%20Renmei/%5BKunagisa%5D%20haibane_renmei_03%5Bh264.vorbis%5D%5Bcestfait%5D.srt",
		},
		{
			Title: "JoJo no Kimyou na Bouken Stone Ocean [Netflix] 01.srt",
			Time:  1638405552,
			Url:   "https://kitsunekko.net/subtitles/japanese/JoJo%20no%20Kimyou%20na%20Bouken:%20Stone%20Ocean/JoJo%20no%20Kimyou%20na%20Bouken%20Stone%20Ocean%20%5BNetflix%5D%2001.srt",
		},
	}
	_okDataStructSubsInfo = []animesubs.SubsInfo{
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
	_brokenDataTime    = "[Judas] Go-Toubun no Hanayome (Season 2).zip|broken|someurl"
	_okDataUrlFilePath = "okurls"
	_okDataUrl         = "[SubsPlease] Kawaii dake ja Nai Shikimori-san - 03 (1080p) [3B413086].mkv|1650737801|https://nyaa.si/view/1518996/torrent\n[SubsPlease] Spy x Family - 03 (1080p) [369CC4DE].mkv|1650727912|https://nyaa.si/view/1518809/torrent\n"
	_okDataStructUrl   = []animeurlservice.AnimeUrlInfo{
		{
			Title:       "[SubsPlease] Kawaii dake ja Nai Shikimori-san - 03 (1080p) [3B413086].mkv",
			TimeUpdated: time.Unix(1650737801, 0),
			Url:         "https://nyaa.si/view/1518996/torrent",
		},
		{
			Title:       "[SubsPlease] Spy x Family - 03 (1080p) [369CC4DE].mkv",
			TimeUpdated: time.Unix(1650727912, 0),
			Url:         "https://nyaa.si/view/1518809/torrent",
		},
	}
)

type mockfileio struct{}

// ReadFile implements fileio.FileIO
func (*mockfileio) ReadFile(filePath string) ([]byte, error) {
	fileMap := map[string]string{
		_okDataSubsFilePath: _okDataSubs,
		_okDataUrlFilePath:  _okDataUrl,
	}

	return []byte(fileMap[filePath]), nil
}

// AppendToFile implements fileio.FileIO
func (*mockfileio) AppendToFile(data []byte, filePath string) error {
	return nil
}

// SaveToFile implements fileio.FileIO
func (*mockfileio) SaveToFile(data []byte, filePath string) error {
	return nil
}

var _ fileio.FileIO = (*mockfileio)(nil)

func (f *mockfileio) SaveResponseToFile(data *http.Response, filePath string) error {
	return nil
}

func TestFormAnimeDataValid(t *testing.T) {
	assert := assert.New(t)

	actual, err := formAnimeData(_okDataSubs)
	assert.NoError(err)

	for i := range _okDataStructSubs {
		assert.True(_okDataStructSubs[i].Equal(actual[i]))
	}
}

func TestFormAnimeDataInvalid(t *testing.T) {
	assert := assert.New(t)

	actual, err := formAnimeData(_brokenDataTime)
	assert.Error(err)
	assert.Len(actual, 0)
}

func TestFormAnimeUrlValid(t *testing.T) {
	assert := assert.New(t)

	actual, err := formAnimeUrlInfos(_okDataUrl)
	assert.NoError(err)
	assert.Len(actual, len(_okDataStructUrl))
	for i := range _okDataStructUrl {
		assert.True(_okDataStructUrl[i].Equal(actual[i]))
	}
}

func TestFormAnimeUrlError(t *testing.T) {
	assert := assert.New(t)

	actual, err := formAnimeUrlInfos(_brokenDataTime)
	assert.Error(err)
	assert.Len(actual, 0)
}

func TestFormAnimeSubsValid(t *testing.T) {
	assert := assert.New(t)
	data := _okDataStructSubsInfo
	actual, err := formAnimeSubs(_okDataSubs)

	assert.NoError(err)
	assert.Len(actual, len(data))
	for i := range data {
		assert.True(data[i].Equal(actual[i]))
	}
}
func TestFormAnimeSubsError(t *testing.T) {
	assert := assert.New(t)

	actual, err := formAnimeSubs(_brokenDataTime)
	assert.Error(err)
	assert.Len(actual, 0)
}

func TestFormFileAnimeDataFromAnimeUrl(t *testing.T) {
	assert := assert.New(t)

	data := animeurlservice.AnimeUrlInfo{
		Title:       "test",
		TimeUpdated: time.Unix(0, 0),
		Url:         "testUrl",
	}

	expected := FileAnimeData{
		Title: data.Title,
		Time:  data.TimeUpdated.Unix(),
		Url:   data.Url,
	}

	actual := formFileAnimeDataFromAnimeUrl(data)

	assert.Equal(expected, actual)
}

func TestFormFileAnimeDataFromSubs(t *testing.T) {
	assert := assert.New(t)

	data := animesubs.SubsInfo{
		Title:       "test",
		TimeUpdated: time.Unix(0, 0),
		Url:         "testUrl",
	}

	expected := FileAnimeData{
		Title: data.Title,
		Time:  data.TimeUpdated.Unix(),
		Url:   data.Url,
	}

	actual := formFileAnimeDataFromAnimeSubs(data)

	assert.Equal(expected, actual)
}

func TestFormStringDataToWrite(t *testing.T) {
	assert := assert.New(t)

	data := FileAnimeData{
		Title: "[Judas] Go-Toubun no Hanayome (Season 2).zip",
		Url:   "https://kitsunekko.net/subtitles/japanese/Gotoubun_No_Hanayome/%5BJudas%5D%20Go-Toubun%20no%20Hanayome%20%28Season%202%29.zip",
		Time:  1643792994,
	}
	expected := "[Judas] Go-Toubun no Hanayome (Season 2).zip|1643792994|https://kitsunekko.net/subtitles/japanese/Gotoubun_No_Hanayome/%5BJudas%5D%20Go-Toubun%20no%20Hanayome%20%28Season%202%29.zip"

	actual := formStringDataToWrite(data)

	assert.Equal(expected, actual)
}

func TestGetAnimeUrlByName(t *testing.T) {
	assert := assert.New(t)

	db := &fileDatabase{
		animeUrlPathFile:  _okDataUrlFilePath,
		animeSubsPathFile: _okDataSubsFilePath,
		storagePath:       "",
		fileIo:            &mockfileio{},
		logger:            &zap.SugaredLogger{},
	}

	data := "[SubsPlease] Kawaii dake ja Nai Shikimori-san - 03 (1080p) [3B413086].mkv"

	expected := animeurlservice.AnimeUrlInfo{
		Title:       "[SubsPlease] Kawaii dake ja Nai Shikimori-san - 03 (1080p) [3B413086].mkv",
		TimeUpdated: time.Unix(1650737801, 0),
		Url:         "https://nyaa.si/view/1518996/torrent",
	}

	actual, err := db.GetAnimeUrlByName(data)
	assert.NoError(err)
	assert.True(expected.Equal(actual))
}

func TestGetAnimeSubsByName(t *testing.T) {
	assert := assert.New(t)

	db := &fileDatabase{
		animeUrlPathFile:  _okDataUrlFilePath,
		animeSubsPathFile: _okDataSubsFilePath,
		storagePath:       "",
		fileIo:            &mockfileio{},
		logger:            &zap.SugaredLogger{},
	}

	data := "[Judas] Go-Toubun no Hanayome (Season 2).zip"

	expected := animesubs.SubsInfo{
		Title:       "[Judas] Go-Toubun no Hanayome (Season 2).zip",
		Url:         "https://kitsunekko.net/subtitles/japanese/Gotoubun_No_Hanayome/%5BJudas%5D%20Go-Toubun%20no%20Hanayome%20%28Season%202%29.zip",
		TimeUpdated: time.Unix(1643792994, 0),
	}

	actual, err := db.GetAnimeSubByName(data)
	assert.NoError(err)
	assert.True(expected.Equal(actual))
}

func TestAddAnimeUrls(t *testing.T) {
	assert := assert.New(t)
	db := &fileDatabase{
		animeUrlPathFile:  _okDataUrlFilePath,
		animeSubsPathFile: _okDataSubsFilePath,
		storagePath:       "",
		fileIo:            &mockfileio{},
		logger:            zap.L().Sugar(),
	}

	err := db.AddAnimeUrls(_okDataStructUrl...)
	assert.NoError(err)
}

func TestAddAnimeUrlsZeroLen(t *testing.T) {
	assert := assert.New(t)
	db := &fileDatabase{
		animeUrlPathFile:  _okDataUrlFilePath,
		animeSubsPathFile: _okDataSubsFilePath,
		storagePath:       "",
		fileIo:            &mockfileio{},
		logger:            zap.L().Sugar(),
	}

	err := db.AddAnimeUrls(nil...)
	assert.NoError(err)
}

func TestAddAnimeSubs(t *testing.T) {
	assert := assert.New(t)
	db := &fileDatabase{
		animeUrlPathFile:  _okDataUrlFilePath,
		animeSubsPathFile: _okDataSubsFilePath,
		storagePath:       "",
		fileIo:            &mockfileio{},
		logger:            zap.L().Sugar(),
	}

	err := db.AddSubs(_okDataStructSubsInfo...)
	assert.NoError(err)
}

func TestAddAnimeSubsZeroLen(t *testing.T) {
	assert := assert.New(t)
	db := &fileDatabase{
		animeUrlPathFile:  _okDataUrlFilePath,
		animeSubsPathFile: _okDataSubsFilePath,
		storagePath:       "",
		fileIo:            &mockfileio{},
		logger:            zap.L().Sugar(),
	}

	err := db.AddSubs(nil...)
	assert.NoError(err)
}
