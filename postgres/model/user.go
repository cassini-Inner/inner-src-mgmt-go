package model

// Users table model
type User struct {
	Id          string
	Email       string
	Name        string
	Role        string
	Department  string
	Bio         string
	PhotoUrl    string `db:"photo_url"`
	Contact     string
	TimeCreated string `db:"time_created"`
	TimeUpdated string `db:"time_updated"`
	IsDeleted   bool   `db:"is_deleted"`
}
