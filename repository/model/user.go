package model

import "database/sql"

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
	TimeCreated string  `db:"time_created"`
	TimeUpdated string  `db:"time_updated"`
	IsDeleted   bool    `db:"is_deleted"`
	GithubUrl   sql.NullString `db:"github_url"`
	Onboarded   bool
	GithubId    sql.NullString `db:"github_id"`
	GithubName  sql.NullString `db:"github_name"`
}
