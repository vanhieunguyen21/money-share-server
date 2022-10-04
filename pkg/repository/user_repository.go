package repository

import (
	"money_share/pkg/model"
)

type UserRepository interface {
	GetById(userId uint) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	CheckUsernameAvailability(username string) (bool, error)
	ValidateUsernameAndUserID(username string, userID uint) (bool, error)
	Create(user *model.User) (*model.User, error)
	Update(user *model.User) (*model.User, error)
	Delete(userId uint) error
}