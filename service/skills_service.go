package service

import "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"

type SkillsService interface {
	GetMatchingSkills(query string, limit *int) ([]*model.GlobalSkill, error)
}
