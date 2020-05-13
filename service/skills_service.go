package service

import "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"

type SkillsService interface {
	GetMatchingSkills(query *string) ([]*model.GlobalSkill, error)
}
