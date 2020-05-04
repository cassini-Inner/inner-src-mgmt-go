package repository

import (
	"database/sql"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
	"github.com/jmoiron/sqlx"
)

type UsersRepo interface {
	Repository
	RemoveUserSkillsByUserId(tx *sqlx.Tx, userID string) error
	GetByIdTx(tx *sqlx.Tx, userId string) (*dbmodel.User, error)
	GetById(userId string) (*dbmodel.User, error)
	GetByEmailId(emailId string) (*dbmodel.User, error)
	GetByGithubId(githubId string) (*dbmodel.User, error)
	CreateNewUser(tx *sqlx.Tx, user *dbmodel.User) (*dbmodel.User, error)
	CountUsersByGithubId(tx *sqlx.Tx, githubId sql.NullString) (int, error)
	UpdateUser(tx *sqlx.Tx, updatedUserInformation *dbmodel.User) (*dbmodel.User, error)
}
