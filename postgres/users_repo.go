package postgres

import (
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/postgres/models"
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
	err := u.db.QueryRowx(query, input.Email).Scan(&user)
	if err != nil {
		var lastInsertId string
		query = "INSERT INTO users (name, email, photo_url) VALUES($1, $2, $3) RETURNING id"
		err = u.db.QueryRowx(query, input.Name, input.Email, input.PhotoURL).Scan(&lastInsertId)
		user.Id = lastInsertId
		user.Email = input.Email
		user.Name = input.Name
		user.PhotoUrl = input.PhotoURL
	}
	return &user, err
}

func (u *UsersRepo) UpdateUser(input *gqlmodel.CreateUserInput) (*gqlmodel.User, error) {
	panic("not implemented")
}

func (u *UsersRepo) GetById(userId string) (*gqlmodel.User, error) {
	panic("not implemented")
}
