package models

// Users table model
type User struct {
	Id          string
	Email       string
	Name        string
	Role        string
	Department  string
	Bio         string
	PhotoUrl    string
	Contact     string
	TimeCreated string
	TimeUpdated string
	IsDeleted   bool
}
