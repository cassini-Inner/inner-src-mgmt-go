package models

import (
	"database/sql"
)

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
