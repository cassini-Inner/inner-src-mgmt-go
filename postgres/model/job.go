package model

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
	IsDeleted   string `db:"is_deleted"`
}
