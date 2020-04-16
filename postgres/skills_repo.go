package postgres

import "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"

type SkillsRepo struct {

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


