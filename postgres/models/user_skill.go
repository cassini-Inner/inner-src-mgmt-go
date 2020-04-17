package models

// UserSkills table model
type UserSkill struct {
	Id          string
	UserId      string
	SkillId     string
	TimeCreated string
	TimeUpdated string
	IsDeleted   bool
}
