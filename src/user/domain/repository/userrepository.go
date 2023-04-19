package repository

import "gobot/src/user/domain/model"

type UserRepository interface {
	GetAll() []*model.User
	Save(*model.User) error
}
