package model

type Milestones struct {
	TotalCounnt *int         `json:"totalCounnt"`
	Milestones  []*Milestone `json:"milestones"`
}
