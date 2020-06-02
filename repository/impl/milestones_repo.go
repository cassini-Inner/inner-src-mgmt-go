package impl

import (
	"context"
	"fmt"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
	"github.com/jmoiron/sqlx"
	"strings"
)

type MilestonesRepoImpl struct {
	db *sqlx.DB
}

func NewMilestonesRepoImpl(db *sqlx.DB) *MilestonesRepoImpl {
	return &MilestonesRepoImpl{db: db}
}

func (m MilestonesRepoImpl) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return m.db.BeginTxx(ctx, nil)
}

func (m MilestonesRepoImpl) CreateMilestones(ctx context.Context, tx *sqlx.Tx, jobId string, milestones []*dbmodel.Milestone) (createdMilestones []*dbmodel.Milestone, err error){
	stmt, valueArgs := m.getInsertMilestonesStatement(milestones, jobId)
	stmt = tx.Rebind(stmt)

	// get the ids of newly inserted milestones
	milestonesInsertResult, err := tx.QueryxContext(ctx, stmt, valueArgs...)
	defer milestonesInsertResult.Close()
	if err != nil {
		return nil, err
	}

	for milestonesInsertResult.Next() {
		var tempMilestone dbmodel.Milestone
		err := milestonesInsertResult.StructScan(&tempMilestone)
		if err != nil {
			return nil, err
		}
		createdMilestones = append(createdMilestones, &tempMilestone)
	}

	return createdMilestones, nil
}

func (m MilestonesRepoImpl) GetByJobId(tx sqlx.Ext, jobId string) ([]*dbmodel.Milestone, error) {
	rows, err := tx.Queryx(selectMilestonesByJobId, jobId)
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

func (m MilestonesRepoImpl) GetIdsByJobId(tx sqlx.Ext, jobId string) (result []string, err error) {
	milestones, err := m.GetByJobId(tx, jobId)
	if err != nil {
		return nil, err
	}

	for _, milestone := range milestones {
		result = append(result, milestone.Id)
	}

	return result, nil
}

func (m MilestonesRepoImpl) GetById(milestoneId string) (*dbmodel.Milestone, error) {
	var milestone dbmodel.Milestone
	err := m.db.QueryRowx(selectMilestoneByIdQuery, milestoneId).StructScan(&milestone)
	if err != nil {
		return nil, err
	}
	return &milestone, nil
}

func (m MilestonesRepoImpl) GetAuthor(milestoneId string) (*dbmodel.User, error) {
	var user dbmodel.User
	err := m.db.QueryRowx(selectUserByMilestoneIdQuery, milestoneId).StructScan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m MilestonesRepoImpl) ForceAutoUpdateMilestoneStatusByJobID(ctx context.Context, tx *sqlx.Tx, jobId string) error {
	_, err := tx.ExecContext(ctx, updateMilestoneStatusByJobIdForce, jobId)
	if err != nil {
		return err
	}
	return nil
}

func (m MilestonesRepoImpl) ForceAutoUpdateMilestoneStatusByMilestoneId(ctx context.Context, tx *sqlx.Tx, milestoneID string) error {
	_, err := tx.ExecContext(ctx, updateMilestoneStatusByMilestoneIDForce, milestoneID)
	if err != nil {
		return err
	}
	return nil
}

func (m MilestonesRepoImpl) MarkMilestonesCompleted(tx *sqlx.Tx, ctx context.Context, milestoneIds ...string) error {
	stmt, args, err := sqlx.In(updateMilestoneStatusCompleted, milestoneIds)
	if err != nil {
		return err
	}

	stmt = tx.Rebind(stmt)
	_, err = tx.ExecContext(ctx, stmt, args...)
	if err != nil {
		return nil
	}
	return nil
}

func (m MilestonesRepoImpl) DeleteMilestonesByJobId(tx *sqlx.Tx, jobID string) error {
	_, err := tx.Exec(deleteMilestonesByJobId, jobID)
	if err != nil {
		return err
	}
	return nil
}


// prepares a statement to insert multiple milestones in a single statement
func (m MilestonesRepoImpl) getInsertMilestonesStatement(milestoneInputs []*dbmodel.Milestone, insertedJobId string) (string, []interface{}) {
	var valueStrings []string
	var valueArgs []interface{}
	for _, milestone := range milestoneInputs {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, milestone.Title, milestone.Description, milestone.Duration, milestone.Status, milestone.Resolution, insertedJobId)
	}
	stmt := fmt.Sprintf("INSERT INTO milestones (title, description, duration, status, resolution, job_id) VALUES %s returning *",
		strings.Join(valueStrings, ", "))
	return stmt, valueArgs
}
