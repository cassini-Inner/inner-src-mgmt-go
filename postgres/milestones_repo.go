package postgres

import (
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/jmoiron/sqlx"
)

type MilestonesRepo struct {
	db *sqlx.DB
}

func NewMilestonesRepo(db *sqlx.DB) *MilestonesRepo {
	return &MilestonesRepo{db: db}
}

//TODO: Refactor this. Should return a list of milestones
func (m *MilestonesRepo) GetByJobId(jobId string) (*model.Milestones, error) {
	panic("not impl")
}
