package model

// Applications table model
type Application struct {
	Id          string
	MilestoneId string `db:"milestone_id"`
	ApplicantId string `db:"applicant_id"`
	Status      string
	Note        string
	TimeCreated string `db:"time_created"`
	TimeUpdated string `db:"time_updated"`
}
