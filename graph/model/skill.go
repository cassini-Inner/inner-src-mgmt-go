package model

import(
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/postgres/models"
)

type Skill struct {
	ID          string `json:"id"`
	CreatedBy   string `json:"createdBy"`
	Value       string `json:"value"`
	CreatedTime string `json:"createdTime"`
}

func (gqlSkill *Skill) mapDbToGql(dbSkill dbmodel.globalskill) {
	if dbSkill.Id != nil {
		gqlSkill.ID = dbSkill.Id
	}

	if dbSkill.CreatedBy != nil {
		gqlSkill.CreatedBy = dbSkill.CreatedBy
	}

	if dbSkill.Value != nil {
		gqlSkill.Value = dbSkill.Value
	}

	if dbSkill.TimeCreated != nil {
		gqlSkill.CreatedTime = dbSkill.TimeCreated
	}
}