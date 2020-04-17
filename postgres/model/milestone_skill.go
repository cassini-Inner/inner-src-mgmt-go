package models

// MilestoneSkills table model
type MilestoneSkill struct {
	Id          string
	MilestoneId string `db:"milestone_id"`
	SkillId     string `db:"skill_id"`
	TimeCreated string `db:"time_created"`
	TimeUpdated string `db:"time_updated"`
	IsDeleted   bool
}
