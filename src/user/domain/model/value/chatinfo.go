package value

import (
	"errors"
	"gobot/src/common/validator"
)

type ChatInfo struct {
	ChatId string
}

func NewChatInfo(chatId string) (*ChatInfo, error) {
	if !validator.ValidateString(chatId) {
		return nil, errors.New("invalid chat id")
	}

	return &ChatInfo{chatId}, nil
}
