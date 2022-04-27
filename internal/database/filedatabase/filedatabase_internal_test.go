package filedatabase

import (
	"encoding/json"
	"fmt"
	"gobot/pkg/fileio"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// import (
// 	"gobot/pkg/animesubs"
// 	"gobot/pkg/animeurlservice"
// 	"gobot/pkg/fileio"
// 	"net/http"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// 	"go.uber.org/zap"
// )

// var (
// 	_okDataSubsFilePath = "oksubs"
// 	_okDataSubs         = "[Judas] Go-Toubun no Hanayome (Season 2).zip|1643792994|https://kitsunekko.net/subtitles/japanese/Gotoubun_No_Hanayome/%5BJudas%5D%20Go-Toubun%20no%20Hanayome%20%28Season%202%29.zip\n[Kunagisa] haibane renmei 03[h264.vorbis][cestfait].srt|1622283917|https://kitsunekko.net/subtitles/japanese/Haibane%20Renmei/%5BKunagisa%5D%20haibane_renmei_03%5Bh264.vorbis%5D%5Bcestfait%5D.srt\nJoJo no Kimyou na Bouken Stone Ocean [Netflix] 01.srt|1638405552|https://kitsunekko.net/subtitles/japanese/JoJo%20no%20Kimyou%20na%20Bouken:%20Stone%20Ocean/JoJo%20no%20Kimyou%20na%20Bouken%20Stone%20Ocean%20%5BNetflix%5D%2001.srt\n"
// 	_okDataStructSubs   = []FileAnimeData{
// 		{
// 			Title: "[Judas] Go-Toubun no Hanayome (Season 2).zip",
// 			Url:   "https://kitsunekko.net/subtitles/japanese/Gotoubun_No_Hanayome/%5BJudas%5D%20Go-Toubun%20no%20Hanayome%20%28Season%202%29.zip",
// 			Time:  1643792994,
// 		},
// 		{
// 			Title: "[Kunagisa] haibane renmei 03[h264.vorbis][cestfait].srt",
// 			Time:  1622283917,
// 			Url:   "https://kitsunekko.net/subtitles/japanese/Haibane%20Renmei/%5BKunagisa%5D%20haibane_renmei_03%5Bh264.vorbis%5D%5Bcestfait%5D.srt",
// 		},
// 		{
// 			Title: "JoJo no Kimyou na Bouken Stone Ocean [Netflix] 01.srt",
// 			Time:  1638405552,
// 			Url:   "https://kitsunekko.net/subtitles/japanese/JoJo%20no%20Kimyou%20na%20Bouken:%20Stone%20Ocean/JoJo%20no%20Kimyou%20na%20Bouken%20Stone%20Ocean%20%5BNetflix%5D%2001.srt",
// 		},
// 	}
// 	_okDataStructSubsInfo = []animesubs.SubsInfo{
// 		{
// 			Title:       "[Judas] Go-Toubun no Hanayome (Season 2).zip",
// 			Url:         "https://kitsunekko.net/subtitles/japanese/Gotoubun_No_Hanayome/%5BJudas%5D%20Go-Toubun%20no%20Hanayome%20%28Season%202%29.zip",
// 			TimeUpdated: time.Unix(1643792994, 0),
// 		},
// 		{
// 			Title:       "[Kunagisa] haibane renmei 03[h264.vorbis][cestfait].srt",
// 			TimeUpdated: time.Unix(1622283917, 0),
// 			Url:         "https://kitsunekko.net/subtitles/japanese/Haibane%20Renmei/%5BKunagisa%5D%20haibane_renmei_03%5Bh264.vorbis%5D%5Bcestfait%5D.srt",
// 		},
// 		{
// 			Title:       "JoJo no Kimyou na Bouken Stone Ocean [Netflix] 01.srt",
// 			TimeUpdated: time.Unix(1638405552, 0),
// 			Url:         "https://kitsunekko.net/subtitles/japanese/JoJo%20no%20Kimyou%20na%20Bouken:%20Stone%20Ocean/JoJo%20no%20Kimyou%20na%20Bouken%20Stone%20Ocean%20%5BNetflix%5D%2001.srt",
// 		},
// 	}
// 	_brokenDataTime    = "[Judas] Go-Toubun no Hanayome (Season 2).zip|broken|someurl"
// 	_okDataUrlFilePath = "okurls"
// 	_okDataUrl         = "[SubsPlease] Kawaii dake ja Nai Shikimori-san - 03 (1080p) [3B413086].mkv|1650737801|https://nyaa.si/view/1518996/torrent\n[SubsPlease] Spy x Family - 03 (1080p) [369CC4DE].mkv|1650727912|https://nyaa.si/view/1518809/torrent\n"
// 	_okDataStructUrl   = []animeurlservice.AnimeUrlInfo{
// 		{
// 			Title:       "[SubsPlease] Kawaii dake ja Nai Shikimori-san - 03 (1080p) [3B413086].mkv",
// 			TimeUpdated: time.Unix(1650737801, 0),
// 			Url:         "https://nyaa.si/view/1518996/torrent",
// 		},
// 		{
// 			Title:       "[SubsPlease] Spy x Family - 03 (1080p) [369CC4DE].mkv",
// 			TimeUpdated: time.Unix(1650727912, 0),
// 			Url:         "https://nyaa.si/view/1518809/torrent",
// 		},
// 	}
// )

type mockFS struct {
	folder string
	files  []mockFile
}

type mockFile struct {
	name     string
	contents string
}

var filesData = []mockFile{
	{
		name:     "kawaii",
		contents: `{"Time":1650737801,"Title":"[SubsPlease] Kawaii dake ja Nai Shikimori-san - 03 (1080p) [3B413086].mkv","Url":"https://nyaa.si/view/1518996/torrent"}`,
	},
	{
		name:     "spy",
		contents: `{"Time":1650727912,"Title":"[SubsPlease] Spy x Family - 03 (1080p) [369CC4DE].mkv","Url":"https://nyaa.si/view/1518809/torrent"}`,
	},
	{
		name:     "empty",
		contents: "",
	},
}

type mockfileio struct{}

var dirsData = []mockFS{
	{
		folder: "1",
		files:  filesData,
	},
}

// GetFilesInDir implements fileio.FileIO
func (*mockfileio) GetFilesInDir(dirpath string) ([]string, error) {
	dirpath = strings.ReplaceAll(dirpath, "/", "")
	for _, dir := range dirsData {
		if dir.folder == dirpath {
			files := make([]string, 0, len(dir.files))
			for _, file := range dir.files {
				files = append(files, file.name)
			}

			return files, nil
		}
	}

	return nil, fmt.Errorf("no data")
}

// AppendToFile implements fileio.FileIO
func (*mockfileio) AppendToFile(data []byte, filePath string) error {
	return nil
}

// CreateDirectory implements fileio.FileIO
func (*mockfileio) CreateDirectory(dirpath string) error {

	dirpath = strings.ReplaceAll(dirpath, "/", "")
	for _, dir := range dirsData {
		if dir.folder == dirpath {
			return nil
		}
	}

	return fmt.Errorf("Error creating folder")
}

// ReadFile implements fileio.FileIO
func (*mockfileio) ReadFile(filePath string) ([]byte, error) {
	filePathSplit := strings.Split(filePath, "/")
	filePath = filePathSplit[len(filePathSplit)-1]
	for _, file := range filesData {
		if file.name == filePath {
			return []byte(file.contents), nil
		}
	}
	return nil, fmt.Errorf("some error")
}

// SaveResponseToFile implements fileio.FileIO
func (*mockfileio) SaveResponseToFile(data *http.Response, filePath string) error {
	return nil
}

// SaveToFile implements fileio.FileIO
func (*mockfileio) SaveToFile(data []byte, filePath string) error {
	return nil
}

var _ fileio.FileIO = (*mockfileio)(nil)

type mockfileioerror struct{}

// AppendToFile implements fileio.FileIO
func (*mockfileioerror) AppendToFile(data []byte, filePath string) error {
	return fmt.Errorf("")
}

// CreateDirectory implements fileio.FileIO
func (*mockfileioerror) CreateDirectory(dirpath string) error {
	return nil
}

// GetFilesInDir implements fileio.FileIO
func (*mockfileioerror) GetFilesInDir(dirpath string) ([]string, error) {
	return nil, fmt.Errorf("")
}

// ReadFile implements fileio.FileIO
func (*mockfileioerror) ReadFile(filePath string) ([]byte, error) {
	return nil, fmt.Errorf("")
}

// SaveResponseToFile implements fileio.FileIO
func (*mockfileioerror) SaveResponseToFile(data *http.Response, filePath string) error {
	return fmt.Errorf("")
}

// SaveToFile implements fileio.FileIO
func (*mockfileioerror) SaveToFile(data []byte, filePath string) error {
	return fmt.Errorf("")
}

var _ fileio.FileIO = (*mockfileioerror)(nil)

func TestReadFileValid(t *testing.T) {
	db := NewFileDatabase("", zap.L().Sugar())
	db.fileIo = &mockfileio{}

	type readFileTestStruct struct {
		name     string
		data     string
		expected string
	}

	var testCases []readFileTestStruct

	for _, file := range filesData {
		testCases = append(testCases, readFileTestStruct{
			name:     file.name,
			data:     file.name,
			expected: file.contents,
		})
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual, err := db.readFile(testCase.data)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestReadFileNoFile(t *testing.T) {
	assert := assert.New(t)
	db := NewFileDatabase("", zap.L().Sugar())
	db.fileIo = &mockfileio{}

	actual, err := db.readFile("none")
	assert.Error(err)
	assert.Equal(actual, "")
}

func TestFindFileNameInFolder(t *testing.T) {
	db := NewFileDatabase("", zap.L().Sugar())
	db.fileIo = &mockfileio{}

	type readFindFileNameInFolderStruct struct {
		name     string
		folder   string
		data     string
		expected string
	}

	var testCases []readFindFileNameInFolderStruct

	for _, dir := range dirsData {
		testCases = append(testCases, readFindFileNameInFolderStruct{
			name:     dir.folder,
			folder:   dir.folder,
			data:     strings.ReplaceAll(dir.files[0].name, "."+DefaultExtension, ""),
			expected: dir.files[0].name,
		})
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual, err := db.findFileNameInFolder(testCase.folder, testCase.data)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestFindFileNameInFolderNoFolder(t *testing.T) {
	assert := assert.New(t)
	db := NewFileDatabase("", zap.L().Sugar())
	db.fileIo = &mockfileio{}

	actual, err := db.findFileNameInFolder("", "")
	assert.Error(err)
	assert.Equal(actual, "")
}

func TestFindFileNameInFolderNotFound(t *testing.T) {
	db := NewFileDatabase("", zap.L().Sugar())
	db.fileIo = &mockfileio{}

	type readFindFileNameInFolderStruct struct {
		name     string
		folder   string
		data     string
		expected string
	}

	var testCases []readFindFileNameInFolderStruct

	for _, dir := range dirsData {
		testCases = append(testCases, readFindFileNameInFolderStruct{
			name:     dir.folder,
			folder:   dir.folder,
			data:     "",
			expected: "",
		})
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual, err := db.findFileNameInFolder(testCase.folder, testCase.data)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestGetFilesInDir(t *testing.T) {
	db := NewFileDatabase("", zap.L().Sugar())
	db.fileIo = &mockfileio{}

	type getFilesInDirStruct struct {
		name     string
		folder   string
		expected []mockFile
	}

	var testCases []getFilesInDirStruct

	for _, dir := range dirsData {
		testCases = append(testCases, getFilesInDirStruct{
			name:     dir.folder,
			folder:   dir.folder,
			expected: dir.files,
		})
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual, err := db.getFilesInDir(testCase.folder)
			assert.NoError(t, err)
			assert.Len(t, actual, len(testCase.expected))
			for i := range actual {
				assert.Equal(t, testCase.expected[i].name, actual[i])
			}
		})
	}
}

func TestGetEntryByNameValid(t *testing.T) {
	db := NewFileDatabase("", zap.L().Sugar())
	db.fileIo = &mockfileio{}
	type getEntryByNameStruct struct {
		name           string
		collectionName string
		data           string
		expected       map[string]interface{}
	}

	var testCases []getEntryByNameStruct

	for _, dir := range dirsData {

		for _, file := range dir.files {
			if file.contents == "" {
				continue
			}
			var result map[string]interface{}
			fmt.Println(file.contents)
			json.Unmarshal([]byte(file.contents), &result)

			testCases = append(testCases, getEntryByNameStruct{
				name:           dir.folder,
				collectionName: dir.folder,
				data:           strings.ReplaceAll(file.name, "."+DefaultExtension, ""),
				expected:       result,
			})
		}

	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual, err := db.GetEntryByName(testCase.collectionName, "", testCase.data)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expected, actual)
		})
	}

}

func TestGetEntryByNameEmpty(t *testing.T) {
	assert := assert.New(t)
	db := NewFileDatabase("", zap.L().Sugar())
	db.fileIo = &mockfileio{}

	var expected map[string]interface{}
	actual, err := db.GetEntryByName(dirsData[0].folder, "", "")
	assert.NoError(err)
	assert.Equal(expected, actual)
}

func TestGetEntryByNameErrorCreatingFolder(t *testing.T) {
	assert := assert.New(t)
	db := NewFileDatabase("", zap.L().Sugar())
	db.fileIo = &mockfileio{}

	var expected map[string]interface{}
	actual, err := db.GetEntryByName("", "", "")
	assert.Error(err)
	assert.Equal(expected, actual)
}
func TestGetEntryByNameErrorGettingFiles(t *testing.T) {
	assert := assert.New(t)
	db := NewFileDatabase("", zap.L().Sugar())
	db.fileIo = &mockfileioerror{}

	var expected map[string]interface{}
	actual, err := db.GetEntryByName("", "", "")
	assert.Error(err)
	assert.Equal(expected, actual)
}

// func (f *mockfileio) SaveResponseToFile(data *http.Response, filePath string) error {
// 	return nil
// }

// func TestFormAnimeDataValid(t *testing.T) {
// 	assert := assert.New(t)

// 	actual, err := formAnimeData(_okDataSubs)
// 	assert.NoError(err)

// 	for i := range _okDataStructSubs {
// 		assert.True(_okDataStructSubs[i].Equal(actual[i]))
// 	}
// }

// func TestFormAnimeDataInvalid(t *testing.T) {
// 	assert := assert.New(t)

// 	actual, err := formAnimeData(_brokenDataTime)
// 	assert.Error(err)
// 	assert.Len(actual, 0)
// }

// func TestFormAnimeUrlValid(t *testing.T) {
// 	assert := assert.New(t)

// 	actual, err := formAnimeUrlInfos(_okDataUrl)
// 	assert.NoError(err)
// 	assert.Len(actual, len(_okDataStructUrl))
// 	for i := range _okDataStructUrl {
// 		assert.True(_okDataStructUrl[i].Equal(actual[i]))
// 	}
// }

// func TestFormAnimeUrlError(t *testing.T) {
// 	assert := assert.New(t)

// 	actual, err := formAnimeUrlInfos(_brokenDataTime)
// 	assert.Error(err)
// 	assert.Len(actual, 0)
// }

// func TestFormAnimeSubsValid(t *testing.T) {
// 	assert := assert.New(t)
// 	data := _okDataStructSubsInfo
// 	actual, err := formAnimeSubs(_okDataSubs)

// 	assert.NoError(err)
// 	assert.Len(actual, len(data))
// 	for i := range data {
// 		assert.True(data[i].Equal(actual[i]))
// 	}
// }
// func TestFormAnimeSubsError(t *testing.T) {
// 	assert := assert.New(t)

// 	actual, err := formAnimeSubs(_brokenDataTime)
// 	assert.Error(err)
// 	assert.Len(actual, 0)
// }

// func TestFormFileAnimeDataFromAnimeUrl(t *testing.T) {
// 	assert := assert.New(t)

// 	data := animeurlservice.AnimeUrlInfo{
// 		Title:       "test",
// 		TimeUpdated: time.Unix(0, 0),
// 		Url:         "testUrl",
// 	}

// 	expected := FileAnimeData{
// 		Title: data.Title,
// 		Time:  data.TimeUpdated.Unix(),
// 		Url:   data.Url,
// 	}

// 	actual := formFileAnimeDataFromAnimeUrl(data)

// 	assert.Equal(expected, actual)
// }

// func TestFormFileAnimeDataFromSubs(t *testing.T) {
// 	assert := assert.New(t)

// 	data := animesubs.SubsInfo{
// 		Title:       "test",
// 		TimeUpdated: time.Unix(0, 0),
// 		Url:         "testUrl",
// 	}

// 	expected := FileAnimeData{
// 		Title: data.Title,
// 		Time:  data.TimeUpdated.Unix(),
// 		Url:   data.Url,
// 	}

// 	actual := formFileAnimeDataFromAnimeSubs(data)

// 	assert.Equal(expected, actual)
// }

// func TestFormStringDataToWrite(t *testing.T) {
// 	assert := assert.New(t)

// 	data := FileAnimeData{
// 		Title: "[Judas] Go-Toubun no Hanayome (Season 2).zip",
// 		Url:   "https://kitsunekko.net/subtitles/japanese/Gotoubun_No_Hanayome/%5BJudas%5D%20Go-Toubun%20no%20Hanayome%20%28Season%202%29.zip",
// 		Time:  1643792994,
// 	}
// 	expected := "[Judas] Go-Toubun no Hanayome (Season 2).zip|1643792994|https://kitsunekko.net/subtitles/japanese/Gotoubun_No_Hanayome/%5BJudas%5D%20Go-Toubun%20no%20Hanayome%20%28Season%202%29.zip"

// 	actual := formStringDataToWrite(data)

// 	assert.Equal(expected, actual)
// }

// func TestGetAnimeUrlByName(t *testing.T) {
// 	assert := assert.New(t)

// 	db := &fileDatabase{
// 		animeUrlPathFile:  _okDataUrlFilePath,
// 		animeSubsPathFile: _okDataSubsFilePath,
// 		storagePath:       "",
// 		fileIo:            &mockfileio{},
// 		logger:            &zap.SugaredLogger{},
// 	}

// 	data := "[SubsPlease] Kawaii dake ja Nai Shikimori-san - 03 (1080p) [3B413086].mkv"

// 	expected := animeurlservice.AnimeUrlInfo{
// 		Title:       "[SubsPlease] Kawaii dake ja Nai Shikimori-san - 03 (1080p) [3B413086].mkv",
// 		TimeUpdated: time.Unix(1650737801, 0),
// 		Url:         "https://nyaa.si/view/1518996/torrent",
// 	}

// 	actual, err := db.GetAnimeUrlByName(data)
// 	assert.NoError(err)
// 	assert.True(expected.Equal(actual))
// }

// func TestGetAnimeSubsByName(t *testing.T) {
// 	assert := assert.New(t)

// 	db := &fileDatabase{
// 		animeUrlPathFile:  _okDataUrlFilePath,
// 		animeSubsPathFile: _okDataSubsFilePath,
// 		storagePath:       "",
// 		fileIo:            &mockfileio{},
// 		logger:            &zap.SugaredLogger{},
// 	}

// 	data := "[Judas] Go-Toubun no Hanayome (Season 2).zip"

// 	expected := animesubs.SubsInfo{
// 		Title:       "[Judas] Go-Toubun no Hanayome (Season 2).zip",
// 		Url:         "https://kitsunekko.net/subtitles/japanese/Gotoubun_No_Hanayome/%5BJudas%5D%20Go-Toubun%20no%20Hanayome%20%28Season%202%29.zip",
// 		TimeUpdated: time.Unix(1643792994, 0),
// 	}

// 	actual, err := db.GetAnimeSubByName(data)
// 	assert.NoError(err)
// 	assert.True(expected.Equal(actual))
// }

// func TestAddAnimeUrls(t *testing.T) {
// 	assert := assert.New(t)
// 	db := &fileDatabase{
// 		animeUrlPathFile:  _okDataUrlFilePath,
// 		animeSubsPathFile: _okDataSubsFilePath,
// 		storagePath:       "",
// 		fileIo:            &mockfileio{},
// 		logger:            zap.L().Sugar(),
// 	}

// 	err := db.AddAnimeUrls(_okDataStructUrl...)
// 	assert.NoError(err)
// }

// func TestAddAnimeUrlsZeroLen(t *testing.T) {
// 	assert := assert.New(t)
// 	db := &fileDatabase{
// 		animeUrlPathFile:  _okDataUrlFilePath,
// 		animeSubsPathFile: _okDataSubsFilePath,
// 		storagePath:       "",
// 		fileIo:            &mockfileio{},
// 		logger:            zap.L().Sugar(),
// 	}

// 	err := db.AddAnimeUrls(nil...)
// 	assert.NoError(err)
// }

// func TestAddAnimeSubs(t *testing.T) {
// 	assert := assert.New(t)
// 	db := &fileDatabase{
// 		animeUrlPathFile:  _okDataUrlFilePath,
// 		animeSubsPathFile: _okDataSubsFilePath,
// 		storagePath:       "",
// 		fileIo:            &mockfileio{},
// 		logger:            zap.L().Sugar(),
// 	}

// 	err := db.AddSubs(_okDataStructSubsInfo...)
// 	assert.NoError(err)
// }

// func TestAddAnimeSubsZeroLen(t *testing.T) {
// 	assert := assert.New(t)
// 	db := &fileDatabase{
// 		animeUrlPathFile:  _okDataUrlFilePath,
// 		animeSubsPathFile: _okDataSubsFilePath,
// 		storagePath:       "",
// 		fileIo:            &mockfileio{},
// 		logger:            zap.L().Sugar(),
// 	}

// 	err := db.AddSubs(nil...)
// 	assert.NoError(err)
// }
