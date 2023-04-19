package model

import (
	"gobot/src/user/domain/model/value"

	"github.com/google/uuid"
)

type User struct {
	Id       uuid.UUID
	UserInfo *value.UserInfo
	ChatInfo *value.ChatInfo
}

func NewTelegramUser(name string, nickname string, password string, chatId string) (*User, error) {
	animeListInfo, err := value.NewAnimeListInfo(nickname, password)
	if err != nil {
		return nil, err
	}

	userInfo, err := value.NewUserInfo(name, animeListInfo)
	if err != nil {
		return nil, err
	}

	chatInfo, err := value.NewChatInfo(chatId)
	if err != nil {
		return nil, err
	}

	id, _ := uuid.NewUUID()

	return &User{id, userInfo, chatInfo}, nil
}
