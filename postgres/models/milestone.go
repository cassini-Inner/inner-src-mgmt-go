package models

// Milestones table model
type Milestone struct {
	Id          string
	AssignedTo  string
	JobId       string
	Title       string
	Description string
	Duration    string
	Resolution  string
	Status      string
	TimeCreated string
	TimeUpdated string
	IsDeleted   bool
}
