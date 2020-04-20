package postgres

import (
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/postgres/model"
	"github.com/jmoiron/sqlx"
)

// TODO: Implement
type UsersRepo struct {
	db *sqlx.DB
}

func NewUsersRepo(db *sqlx.DB) *UsersRepo {
	return &UsersRepo{db: db}
}

//TODO: Refactor. Should return DB model
func (u *UsersRepo) CreateUser(input *gqlmodel.CreateUserInput) (*gqlmodel.User, error) {
	panic("not implemented")
}

//TODO: Implement
func (u *UsersRepo) UpdateUser(input *gqlmodel.CreateUserInput) (*gqlmodel.User, error) {
	panic("not implemented")
}

func (u *UsersRepo) GetById(userId string) (*dbmodel.User, error) {
	var user dbmodel.User
	err := u.db.QueryRowx(getUserByIdQuery, userId).StructScan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

const (
	getUserByIdQuery = `select * from users where users.id = $1`
)