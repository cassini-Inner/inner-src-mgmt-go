package postgres

import (
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/jinzhu/gorm"
)

// TODO: Implement
type UsersRepo struct {
	db *gorm.DB
}

func NewUsersRepo(db *gorm.DB) *UsersRepo {
	return &UsersRepo{db: db}
}

func (u *UsersRepo) CreateUser(input *model.CreateUserInput) (*model.User, error) {
	panic("not implemented")
}

func (u *UsersRepo) UpdateUser(input *model.CreateUserInput) (*model.User, error) {
	panic("not implemented")
}

func (u *UsersRepo) GetById(userId string) (*model.User, error) {
	panic("not implemented")
}
