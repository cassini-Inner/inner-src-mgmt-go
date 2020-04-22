package model

import (
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/postgres/model"
)

type User struct {
	ID          string     `json:"id"`
	Onboarded   bool       `json:"onboarded"`
	Email       string     `json:"email"`
	Name        string     `json:"name"`
	Role        string     `json:"role"`
	Department  string     `json:"department"`
	PhotoURL    string     `json:"photoUrl"`
	GithubURL   string     `json:"githubUrl"`
	Bio         *string    `json:"bio"`
	Contact     *string    `json:"contact"`
	Skills      []*Skill   `json:"skills"`
	TimeCreated string     `json:"timeCreated"`
	TimeUpdated string     `json:"timeUpdated"`
	CreatedJobs []*Job     `json:"createdJobs"`
	AppliedJobs []*Job     `json:"appliedJobs"`
	JobStats    *UserStats `json:"jobStats"`
}

func (u *User) MapDbToGql(dbUser dbmodel.User) {
	u.ID = dbUser.Id
	u.Email = dbUser.Email
	u.Name = dbUser.Name
	u.Onboarded = dbUser.Onboarded
	if dbUser.Role.Valid {
		u.Role = dbUser.Role.String
	}
	if dbUser.Department.Valid {
		u.Department = dbUser.Department.String
	}
	u.Bio = &dbUser.Bio
	if dbUser.Contact.Valid {
		u.Contact = &dbUser.Contact.String
	}
	u.GithubURL = dbUser.GithubUrl
	u.PhotoURL = dbUser.PhotoUrl
	u.TimeCreated = dbUser.TimeCreated
	u.TimeUpdated = dbUser.TimeUpdated
}
