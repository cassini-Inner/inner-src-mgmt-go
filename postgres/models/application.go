package models

// Applications table model
type Application struct {
	Id          string
	MilestoneId string
	ApplicantId string
	Status      string
	Notes       string
	TimeCreated string
	TimeUpdated string
}
