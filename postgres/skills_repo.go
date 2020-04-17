package postgres

import (
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/jmoiron/sqlx"
)

type SkillsRepo struct {
	db *sqlx.DB
}

func NewSkillsRepo(db *sqlx.DB) *SkillsRepo {
	return &SkillsRepo{db: db}
}

func (s *SkillsRepo) GetByJobId(jobId string) ([]*model.Skill, error) {
	panic("not impl")
}

func (s *SkillsRepo) GetByUserId(userId string) ([]*model.Skill, error) {
	panic("not impl")
}

func (s *SkillsRepo) GetByMilestoneId(milestoneId string) ([]*model.Skill, error) {
	panic("not impl")
}
