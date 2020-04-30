package postgres

import (
	"database/sql"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/postgres/model"
	"github.com/jmoiron/sqlx"
)

type UsersRepo interface {
	RemoveUserSkillsByUserId(userID string, tx *sqlx.Tx) error
	GetByIdTx(userId string, tx *sqlx.Tx) (*dbmodel.User, error)
	GetById(userId string) (*dbmodel.User, error)
	GetByEmailId(emailId string) (*dbmodel.User, error)
	GetByGithubId(githubId string) (*dbmodel.User, error)
	CreateNewUser(user *dbmodel.User, tx *sqlx.Tx) (*dbmodel.User, error)
	CountUsersByGithubId(githubId sql.NullString, tx *sqlx.Tx) (int, error)
	UpdateUser(currentUserInfo *dbmodel.User, input *gqlmodel.UpdateUserInput, tx *sqlx.Tx) (*dbmodel.User, error)
}
