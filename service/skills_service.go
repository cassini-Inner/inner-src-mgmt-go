package service

import (
	"github.com/cassini-Inner/inner-src-mgmt-go/repository"
	"github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
)

type SkillsService struct {
	skillsRepo repository.SkillsRepo
}

func NewSkillsService(skillsRepo repository.SkillsRepo) *SkillsService {
	return &SkillsService{skillsRepo: skillsRepo}
}

func (s *SkillsService) GetMatchingSkills(query *string) ([]*model.GlobalSkill, error) {
	if query == nil || *query == "" {
		return s.skillsRepo.GetAll()
	}
	return s.skillsRepo.GetMatchingSkills(query)
}
