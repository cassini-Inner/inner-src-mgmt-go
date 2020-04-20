package model

// GlobalSkills table model
type GlobalSkill struct {
	Id          string
	CreatedBy   string `db:"created_by"`
	Value       string
	TimeCreated string `db:"time_created"`
}
