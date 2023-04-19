package model_test

import (
	"gobot/src/user/domain/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserOk(t *testing.T) {
	name := "test"
	nickname := "nickname"
	password := "password"
	chatId := "chatId"

	_, err := model.NewTelegramUser(name, nickname, password, chatId)
	assert.Nil(t, err)
}

func TestFailAnimeListInfoShouldError(t *testing.T) {
	name := "test"
	nickname := ""
	password := "password"
	chatId := "chatId"

	_, err := model.NewTelegramUser(name, nickname, password, chatId)
	assert.NotNil(t, err)
}

func TestFailUserInfoShouldError(t *testing.T) {
	name := ""
	nickname := "nickname"
	password := "password"
	chatId := "chatId"

	_, err := model.NewTelegramUser(name, nickname, password, chatId)
	assert.NotNil(t, err)
}

func TestFailChatInfoShouldError(t *testing.T) {
	name := "name"
	nickname := "nickname"
	password := "password"
	chatId := ""

	_, err := model.NewTelegramUser(name, nickname, password, chatId)
	assert.NotNil(t, err)
}
