package model

type Comment struct {
	ID          string `json:"id"`
	JobID       string `json:"job"`
	TimeCreated string `json:"timeCreated"`
	TimeUpdated string `json:"timeUpdated"`
	Content     string `json:"content"`
	IsDeleted   bool   `json:"isDeleted"`
	CreatedBy   string `json:"createdBy"`
}
