package model

import (
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/postgres/model"
)

type Skill struct {
	ID          string `json:"id"`
	CreatedBy   string `json:"createdBy"`
	Value       string `json:"value"`
	CreatedTime string `json:"createdTime"`
}

func (s *Skill) MapDbToGql(dbSkill dbmodel.GlobalSkill) {
	s.ID = dbSkill.Id
	s.CreatedBy = dbSkill.CreatedBy
	s.Value = dbSkill.Value
	s.CreatedTime = dbSkill.TimeCreated
}
