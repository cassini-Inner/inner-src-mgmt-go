package model

type UserJobApplication struct {
	ApplicationStatus ApplicationStatus `json:"applicationStatus"`
	UserJobStatus     JobStatus         `json:"userJobStatus"`
	Job               *Job              `json:"job"`
}
