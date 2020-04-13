package model

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
