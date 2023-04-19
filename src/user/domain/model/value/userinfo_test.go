package value_test

import (
	"gobot/src/user/domain/model/value"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyNameShouldError(t *testing.T) {
	name := ""
	animeListInfo := &value.AnimeListInfo{}

	_, err := value.NewUserInfo(name, animeListInfo)
	assert.NotNil(t, err)
}

func TestUserInfoOk(t *testing.T) {
	name := "some"
	animeListInfo := &value.AnimeListInfo{}

	_, err := value.NewUserInfo(name, animeListInfo)
	assert.Nil(t, err)
}
