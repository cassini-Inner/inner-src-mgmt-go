package service

import (
	"github.com/cassini-Inner/inner-src-mgmt-go/repository"
	"github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
	"github.com/jmoiron/sqlx"
)

type SkillsService struct {
	db         *sqlx.DB
	skillsRepo repository.SkillsRepo
}

func NewSkillsService(db *sqlx.DB, skillsRepo repository.SkillsRepo) *SkillsService {
	return &SkillsService{db: db, skillsRepo: skillsRepo}
}

func (s *SkillsService) GetMatchingSkills(query *string) ([]*model.GlobalSkill, error) {
	if query == nil || *query == "" {
		return s.skillsRepo.GetAll()
	}
	return s.skillsRepo.GetMatchingSkills(query)
}
