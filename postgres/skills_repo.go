package postgres

import (
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/jinzhu/gorm"
)

type SkillsRepo struct {
	db *gorm.DB
}

func NewSkillsRepo(db *gorm.DB) *SkillsRepo {
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


