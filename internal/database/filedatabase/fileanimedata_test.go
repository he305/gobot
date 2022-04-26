package filedatabase_test

import (
	"gobot/internal/database/filedatabase"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubsInfoEqual(t *testing.T) {
	assert := assert.New(t)
	expected := true
	first := filedatabase.FileAnimeData{
		Title: "test",
		Time:  0,
		Url:   "http://example.com",
	}
	second := filedatabase.FileAnimeData{
		Title: "test",
		Time:  0,
		Url:   "http://example.com",
	}

	actual := first.Equal(second)
	assert.Equal(expected, actual)
}
