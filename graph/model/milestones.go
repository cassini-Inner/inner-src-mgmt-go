package model

type Milestones struct {
	TotalCount *int         `json:"totalCount"`
	Milestones []*Milestone `json:"milestones"`
}
