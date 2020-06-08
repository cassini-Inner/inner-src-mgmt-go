package repository

import (
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
	"github.com/jmoiron/sqlx"
)

type ReviewsRepo interface {
	Repository
	GetById(id string) (*dbmodel.Review, error)
	GetByIdTx(tx *sqlx.Tx, id string) (*dbmodel.Review, error)
	GetForUserIdMilestoneId(milestoneId, userId string) (*dbmodel.Review, error)
	GetForUserId(userId string) ([]*dbmodel.Review, error)

	Add(tx *sqlx.Tx, review dbmodel.Review) (*dbmodel.Review, error)
	Update(tx *sqlx.Tx, id string, review dbmodel.Review) (*dbmodel.Review, error)
	Delete(tx *sqlx.Tx, id string) (*dbmodel.Review, error)
}
