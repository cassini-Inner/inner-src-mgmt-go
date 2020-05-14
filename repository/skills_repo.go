package repository

import (
	"context"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
	"github.com/jmoiron/sqlx"
)

type SkillsRepo interface {
	Repository
	GetByJobId(jobId string) ([]*dbmodel.GlobalSkill, error)
	GetByUserId(userId string) ([]*dbmodel.GlobalSkill, error)
	GetByMilestoneId(milestoneId string) ([]*dbmodel.GlobalSkill, error)
	GetAll() ([]*dbmodel.GlobalSkill, error)
	FindOrCreateSkills(ctx context.Context, tx *sqlx.Tx, skillsList []string, userId string) (skillsMap map[string]*dbmodel.GlobalSkill, err error)
	MapSkillsToMilestones(ctx context.Context, tx *sqlx.Tx, skillsMap map[string]*dbmodel.GlobalSkill, input *gqlmodel.CreateJobInput, insertedMilestones []*dbmodel.Milestone) (err error)
	AddSkillsToUserSkills(tx *sqlx.Tx, skills map[string]*dbmodel.GlobalSkill, userId string) error
	GetMatchingSkills(query string, limit *int) ([]*dbmodel.GlobalSkill, error)
}
