package models

// Jobs table model
type Job struct {
	Id          string
	CreatedBy   string
	Title       string
	Description string
	Difficulty  string
	Status      string
	TimeCreated string
	TimeUpdated string
	IsDeleted   string
}
