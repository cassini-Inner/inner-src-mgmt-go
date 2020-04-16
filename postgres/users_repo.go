package postgres

import "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"

// TODO: Implement
type UsersRepo struct {
	skillsRepo *SkillsRepo
}

func NewUsersRepo(skillRepo *SkillsRepo) *UsersRepo {
	return &UsersRepo{skillsRepo: skillRepo}
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
