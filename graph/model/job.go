package model

import(
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/postgres/model"
	"strings"
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


func (j *Job) mapDbToGql(dbJob dbmodel.Job) {
	j.ID = dbJob.Id
	j.Title = dbJob.Title
	j.CreatedBy = dbJob.CreatedBy
	j.Desc = dbJob.Description
	j.Difficulty = Difficulty(dbJob.Difficulty)
	j.Status = JobStatus(strings.ToUpper(dbJob.Status))
	j.TimeCreated = dbJob.TimeCreated
	j.TimeUpdated = dbJob.TimeUpdated
}