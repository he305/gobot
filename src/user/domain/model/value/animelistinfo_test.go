package value_test

import (
	"gobot/src/user/domain/model/value"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyNicknameShouldError(t *testing.T) {
	nickname := ""
	pass := "some"

	_, err := value.NewAnimeListInfo(nickname, pass)
	assert.NotNil(t, err)
}

func TestEmptyPasswordShouldError(t *testing.T) {
	nickname := "some"
	pass := ""

	_, err := value.NewAnimeListInfo(nickname, pass)
	assert.NotNil(t, err)
}

func TestAnimeListInfoOk(t *testing.T) {
	nickname := "some"
	pass := "some"

	_, err := value.NewAnimeListInfo(nickname, pass)
	assert.Nil(t, err)
}
