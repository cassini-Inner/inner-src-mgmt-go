package postgres

import (
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/postgres/model"
	"github.com/jmoiron/sqlx"
)

type MilestonesRepo struct {
	db *sqlx.DB
}

func NewMilestonesRepo(db *sqlx.DB) *MilestonesRepo {
	return &MilestonesRepo{db: db}
}

func (m *MilestonesRepo) GetByJobId(jobId string) ([]*dbmodel.Milestone, error) {
	rows, err := m.db.Queryx(selectMilestonesByJobId, jobId)
	if err != nil {
		return nil, err
	}

	var milestones []*dbmodel.Milestone
	for rows.Next() {
	var milestone dbmodel.Milestone
		rows.StructScan(&milestone)
		milestones = append(milestones, &milestone)
	}
	return milestones, nil
}

func (m *MilestonesRepo) GetById(milestoneId string) (*dbmodel.Milestone, error) {
	var milestone dbmodel.Milestone
	err :=m.db.QueryRowx(selectMilestoneByIdQuery, milestoneId).StructScan(&milestone)
	if err != nil {
		return nil, err
	}
	return &milestone, nil
}

const (
	selectMilestonesByJobId = `SELECT * FROM milestones WHERE job_id = $1 and is_deleted = false`
	selectMilestoneByIdQuery = `SELECT * FROM milestones WHERE id = $1 and is_deleted=false`
)