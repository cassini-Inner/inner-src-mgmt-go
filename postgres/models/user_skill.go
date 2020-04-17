package models

// UserSkills table model
type UserSkill struct {
	Id          string
	UserId      string `db:"user_id"`
	SkillId     string `db:"skill_id"`
	TimeCreated string `db:"time_created"`
	TimeUpdated string `db:"time_updated"`
	IsDeleted   bool   `db:"is_deleted"`
}
