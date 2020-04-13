package model

type Skill struct {
	ID          string `json:"id"`
	CreatedBy   string `json:"createdBy"`
	Value       string `json:"value"`
	CreatedTime string `json:"createdTime"`
}
