package repository

import "FoodStore-AdvProg2/domain"

type UserRepository interface {
	Save(user domain.User) (string, error)
	FindByUsername(username string) (domain.User, error)
	FindByID(id string) (domain.User, error)
	SaveToken(token domain.Token) error
	FindUserIDByToken(token string) (string, error)
}