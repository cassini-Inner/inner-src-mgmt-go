package model

import "database/sql"

// Applications table model
type Application struct {
	Id          string
	MilestoneId string `db:"milestone_id"`
	ApplicantId string `db:"applicant_id"`
	Status      string
	Note        sql.NullString
	TimeCreated string `db:"time_created"`
	TimeUpdated string `db:"time_updated"`
}

// Users table model
type User struct {
	Id          string
	Email       sql.NullString
	Name        sql.NullString
	Role        sql.NullString
	Department  sql.NullString
	Bio         sql.NullString
	PhotoUrl    sql.NullString `db:"photo_url"`
	Contact     sql.NullString
	TimeCreated string         `db:"time_created"`
	TimeUpdated string         `db:"time_updated"`
	IsDeleted   bool           `db:"is_deleted"`
	GithubUrl   sql.NullString `db:"github_url"`
	Onboarded   bool
	GithubId    sql.NullString `db:"github_id"`
	GithubName  sql.NullString `db:"github_name"`
}

// UserSkills table model
type UserSkill struct {
	Id          string
	UserId      string `db:"user_id"`
	SkillId     string `db:"skill_id"`
	TimeCreated string `db:"time_created"`
	TimeUpdated string `db:"time_updated"`
	IsDeleted   bool   `db:"is_deleted"`
}

// MilestoneSkills table model
type MilestoneSkill struct {
	Id          string
	MilestoneId string `db:"milestone_id"`
	SkillId     string `db:"skill_id"`
	TimeCreated string `db:"time_created"`
	TimeUpdated string `db:"time_updated"`
	IsDeleted   bool
}

// Milestones table model
type Milestone struct {
	Id          string
	AssignedTo  sql.NullString `db:"assigned_to"`
	JobId       string         `db:"job_id"`
	Title       string
	Description string
	Duration    string
	Resolution  string
	Status      string
	TimeCreated string `db:"time_created"`
	TimeUpdated string `db:"time_updated"`
	IsDeleted   bool   `db:"is_deleted"`
}

// Jobs table model
type Job struct {
	Id          string
	CreatedBy   string `db:"created_by"`
	Title       string
	Description string
	Difficulty  string
	Status      string
	TimeCreated string `db:"time_created"`
	TimeUpdated string `db:"time_updated"`
	IsDeleted   bool   `db:"is_deleted"`
}

// GlobalSkills table model
type GlobalSkill struct {
	Id          string
	CreatedBy   string `db:"created_by"`
	Value       string
	TimeCreated string `db:"time_created"`
}

// Discussions table model
type Discussion struct {
	Id          string
	JobId       string `db:"job_id"`
	CreatedBy   string `db:"created_by"`
	Content     string
	TimeCreated string `db:"time_created"`
	TimeUpdated string `db:"time_updated"`
	IsDeleted   bool   `db:"is_deleted"`
}

type Review struct {
	Id string
	Rating int
	Remark string
	MilestoneId string `db:"milestone_id"`
	UserId string `db:"user_id"`
	TimeCreated string `db:"time_created"`
	TimeUpdated string `db:"time_updated"`
	IsDeleted   bool   `db:"is_deleted"`
}
