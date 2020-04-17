package model

import(
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/postgres/models"
)

type User struct {
	ID          string     `json:"id"`
	Email       string     `json:"email"`
	Name        string     `json:"name"`
	Role        string     `json:"role"`
	Department  string     `json:"department"`
	PhotoURL    string     `json:"photoUrl"`
	Bio         *string    `json:"bio"`
	Contact     *string    `json:"contact"`
	Skills      []*Skill   `json:"skills"`
	TimeCreated string     `json:"timeCreated"`
	TimeUpdated string     `json:"timeUpdated"`
	CreatedJobs []*Job     `json:"createdJobs"`
	AppliedJobs []*Job     `json:"appliedJobs"`
	JobStats    *UserStats `json:"jobStats"`
}

func (u *User) mapDbToGql(dbUser dbmodel.User) {
		u.ID = dbUser.Id
		u.Email = dbUser.Email
		u.Name = dbUser.Name
		u.Role = dbUser.Role
		u.Department = dbUser.Department
		u.PhotoURL = dbUser.PhotoUrl
		u.Bio = &dbUser.Bio
		u.Contact = &dbUser.Contact
		u.TimeCreated = dbUser.TimeCreated
		u.TimeUpdated = dbUser.TimeUpdated
}