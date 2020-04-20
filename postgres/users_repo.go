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

func (u *UsersRepo) CreateUser(input *gqlmodel.CreateUserInput) (*dbmodel.User, error) {
	var user dbmodel.User
	query := "SELECT * FROM users WHERE email = $1 AND is_deleted = FALSE"
	err := u.db.QueryRowx(query, input.Email).StructScan(&user)

	if err != nil {
		var lastInsertId string
		query = "INSERT INTO users (name, email, photo_url) VALUES($1, $2, $3) RETURNING id"
		err = u.db.QueryRowx(query, input.Name, input.Email, input.PhotoURL).Scan(&lastInsertId)
		if err != nil {
			return nil, err
		}
		user.Id = lastInsertId
		user.Email = input.Email
		user.Name = input.Name
		user.PhotoUrl = input.PhotoURL
	}
	return &user, err
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
	getUserByIdQuery = `select * from users where users.id = $1 and users.is_deleted = false`
)