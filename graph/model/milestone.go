package model

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
