package model

import "database/sql"

// Users table model
type User struct {
	Id          string
	Onboarded	bool
	Email       string
	Name        string
	Role        sql.NullString
	Department  sql.NullString
	Bio         string
	PhotoUrl    string `db:"photo_url"`
	GithubUrl   string `db:"github_url"`
	Contact     sql.NullString
	TimeCreated string `db:"time_created"`
	TimeUpdated string `db:"time_updated"`
	IsDeleted   bool   `db:"is_deleted"`
}

