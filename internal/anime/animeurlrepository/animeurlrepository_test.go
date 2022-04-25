package animeurlrepository

import (
	"gobot/internal/database"
	"gobot/pkg/animesubs"
	"gobot/pkg/animeurlservice"
	"testing"

	"github.com/stretchr/testify/assert"
)

type _mockDb struct {
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
	return animeurlservice.AnimeUrlInfo{Title: name}, nil
}

var _ database.Database = (*_mockDb)(nil)

func TestGetAnimeUrlByName(t *testing.T) {
	assert := assert.New(t)

	data := "testName"
	repo := NewAnimeUrlRepository(&_mockDb{})
	expected := animeurlservice.AnimeUrlInfo{Title: data}

	actual, err := repo.GetAnimeUrlByName(data)
	assert.NoError(err)
	assert.Equal(expected.Title, actual.Title)
}

func TestAddAnimeUrls(t *testing.T) {
	assert := assert.New(t)

	repo := NewAnimeUrlRepository(&_mockDb{})
	err := repo.AddAnimeUrls([]animeurlservice.AnimeUrlInfo{}...)
	assert.NoError(err)
}
