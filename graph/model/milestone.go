package model

import(
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/postgres/models"
)

type Milestone struct {
	ID          string     `json:"id"`
	JobID       string     `json:"job"`
	Title       string     `json:"title"`
	TimeCreated string     `json:"timeCreated"`
	TimeUpdated string     `json:"timeUpdated"`
	Desc        string     `json:"desc"`
	Resolution  string     `json:"resolution"`
	Duration    string     `json:"duration"`
	Status      *JobStatus `json:"status"`
	AssignedTo  string     `json:"assignedTo"`
	Skills      []*Skill   `json:"skills"`
}

func (gqlMilestone *Milestone) mapDbToGql(dbMilestone dbmodel.milestone) {

	if dbMilestone.Id != nil {
		gqlMilestone.ID = dbMilestone.Id
	}

	if dbMilestone.JobID != nil {
		gqlMilestone.JobID = dbMilestone.JobID
	}

	if dbMilestone.Title != nil {
		gqlMilestone.Title = dbMilestone.Title
	}

	if dbMilestone.TimeCreated != nil {
		gqlMilestone.TimeCreated = dbMilestone.TimeCreated
	}

	if dbMilestone.TimeUpdated != nil {
		gqlMilestone.TimeUpdated = dbMilestone.TimeUpdated
	}

	if dbMilestone.Description != nil {
		gqlMilestone.Desc = dbMilestone.Description
	}

	if dbMilestone.Resolution != nil {
		gqlMilestone.Resolution = dbMilestone.Resolution
	}

	if dbMilestone.Duration != nil {
		gqlMilestone.Duration = dbMilestone.Duration
	}

	if dbMilestone.Status != nil {
		gqlMilestone.Status = dbMilestone.Status
	}

	if dbMilestone.AssignedTo != nil {
		gqlMilestone.AssignedTo = dbMilestone.AssignedTo
	}
}