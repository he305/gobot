package value

import (
	"errors"
	"gobot/src/common/validator"
)

type AnimeListInfo struct {
	Nickname string
	Password string
}

func NewAnimeListInfo(nickname string, password string) (*AnimeListInfo, error) {
	if !validator.ValidateString(nickname) || !validator.ValidateString(password) {
		return nil, errors.New("invalid anime list info data")
	}

	return &AnimeListInfo{nickname, password}, nil
}
