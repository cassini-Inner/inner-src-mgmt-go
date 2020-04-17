package postgres

import (
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/jinzhu/gorm"
)

type MilestonesRepo struct {
	db *gorm.DB
}

func NewMilestonesRepo(db *gorm.DB) *MilestonesRepo {
	return &MilestonesRepo{db: db}
}

//TODO: Implement
func (m *MilestonesRepo) GetByJobId(jobId string) (*model.Milestones, error) {
	panic("not impl")
}
