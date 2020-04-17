package model

// Milestones table model
type Milestone struct {
	Id          string
	AssignedTo string `db:"assigned_to"`
	JobId       string `db:"job_id"`
	Title       string
	Description string
	Duration    string
	Resolution  string
	Status      string
	TimeCreated string `db:"time_created"`
	TimeUpdated string `db:"time_updated"`
	IsDeleted   bool   `db:"is_deleted"`
}
