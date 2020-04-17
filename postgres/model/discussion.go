package models

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
