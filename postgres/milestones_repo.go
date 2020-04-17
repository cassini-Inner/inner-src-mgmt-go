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
	var milestone dbmodel.Milestone
	var milestones []*dbmodel.Milestone
	query := "SELECT * FROM milestones WHERE job_id = $1"
	rows, err := m.db.Queryx(query, jobId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		rows.StructScan(&milestone)
		milestones = append(milestones, &milestone)
	}
	return milestones, nil
}
