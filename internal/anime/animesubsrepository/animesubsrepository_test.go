package animesubsrepository

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
	return animeurlservice.AnimeUrlInfo{}, nil
}

var _ database.Database = (*_mockDb)(nil)

func TestGetAnimeSubsByName(t *testing.T) {
	assert := assert.New(t)

	data := "testName"
	repo := NewAnimeSubsRepository(&_mockDb{})
	expected := animesubs.SubsInfo{Title: data}

	actual, err := repo.GetAnimeSubsByName(data)
	assert.NoError(err)
	assert.Equal(expected.Title, actual.Title)
}

func TestAddAnimeSubs(t *testing.T) {
	assert := assert.New(t)

	repo := NewAnimeSubsRepository(&_mockDb{})
	err := repo.AddAnimeSubs([]animesubs.SubsInfo{}...)
	assert.NoError(err)
}
