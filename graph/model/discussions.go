package model

type Discussions struct {
	TotalCount  *int       `json:"totalCount"`
	Discussions []*Comment `json:"discussions"`
}
