package models

// Discussions table model
type Discussion struct {
	Id          string
	JobId       string
	CreatedBy   string
	Content     string
	TimeCreated string
	TimeUpdated string
	IsDeleted   bool
}
