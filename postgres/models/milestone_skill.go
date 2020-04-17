package models

// MilestoneSkills table model
type MilestoneSkill struct {
	Id          string
	MilestoneId string
	SkillId     string
	TimeCreated string
	TimeUpdated string
	IsDeleted   bool
}
