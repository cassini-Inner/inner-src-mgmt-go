package model

import(
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/postgres/models"
)

type Job struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	CreatedBy    string        `json:"createdBy"`
	Desc         string        `json:"desc"`
	Duration     string        `json:"duration"`
	Difficulty   Difficulty    `json:"difficulty"`
	Status       JobStatus     `json:"status"`
	TimeCreated  string        `json:"timeCreated"`
	TimeUpdated  string        `json:"timeUpdated"`
	Discussion   *Discussions  `json:"discussion"`
	Milestones   *Milestones   `json:"milestones"`
	Applications *Applications `json:"applications"`
}


func (gqlJob *Job) mapDbToGql(dbJob dbmodel.job) {

	if dbJob.Id != nil {
		gqlJob.ID = dbJob.Id
	}

	if dbJob.Title != nil {
		gqlJob.Title = dbJob.Title
	}

	if dbJob.CreatedBy != nil {
		gqlJob.CreatedBy = dbJob.CreatedBy
	}

	if dbJob.Description != nil {
		gqlJob.Desc = dbJob.Description
	}

	if dbJob.Difficulty != nil {
		gqlJob.Difficulty = dbJob.Difficulty
	}

	if dbJob.Status != nil {
		gqlJob.Status = dbJob.Status
	}

	if dbJob.TimeCreated != nil {
		gqlJob.TimeCreated = dbJob.TimeCreated
	}

	if dbJob.TimeUpdated != nil {
		gqlJob.TimeUpdated = dbJob.TimeUpdated
	}

}