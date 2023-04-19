package value

import (
	"errors"
	"gobot/src/common/validator"
)

type UserInfo struct {
	Name          string
	AnimeListInfo *AnimeListInfo
}

func NewUserInfo(name string, animeListInfo *AnimeListInfo) (*UserInfo, error) {
	if !validator.ValidateString(name) {
		return nil, errors.New("invalid name")
	}

	return &UserInfo{name, animeListInfo}, nil
}
