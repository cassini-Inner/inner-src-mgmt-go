package postgres

import (
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/go-pg/pg/v9"
)

type SkillsRepo struct {
	db *pg.DB
}

func NewSkillsRepo(db *pg.DB) *SkillsRepo {
	return &SkillsRepo{db: db}
}

func (s * SkillsRepo) GetByJobId(jobId string) ([]*model.Skill, error) {
	panic("not impl")
}

func (s * SkillsRepo) GetByUserId(userId string) ([]*model.Skill, error) {
	panic("not impl")
}

func (s * SkillsRepo) GetByMilestoneId(milestoneId string) ([]*model.Skill, error) {
	panic("not impl")
}


