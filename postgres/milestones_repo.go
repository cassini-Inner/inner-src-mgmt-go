package postgres

import (
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/go-pg/pg/v9"
)

type MilestonesRepo struct {
	db *pg.DB
}

func NewMilestonesRepo(db *pg.DB) *MilestonesRepo {
	return &MilestonesRepo{db: db}
}

//TODO: Implement
func (m *MilestonesRepo) GetByJobId(jobId string) (*model.Milestones, error) {
	panic("not impl")
}
