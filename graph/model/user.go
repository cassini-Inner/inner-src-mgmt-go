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

func (gqlUser *User) mapDbToGql(dbUser dbmodel.user) {
	if dbUser.Id != nil {
		gqlUser.ID = dbUser.Id
	}

	if dbUser.Email != nil {
		gqlUser.Email = dbUser.Email
	}

	if dbUser.Name != nil {
		gqlUser.Name = dbUser.Name
	}

	if dbUser.Role != nil {
		gqlUser.Role = dbUser.Role
	}

	if dbUser.Department != nil {
		gqlUser.Department = dbUser.Department
	}

	if dbUser.PhotoURL != nil {
		gqlUser.PhotoURL = dbUser.PhotoURL
	}
	
	if dbUser.Bio != nil {
		gqlUser.Bio = dbUser.Bio
	}

	if dbUser.Contact != nil {
		gqlUser.Contact = dbUser.Contact
	}

	if dbUser.TimeCreated != nil {
		gqlUser.TimeCreated = dbUser.TimeCreated
	}

	if dbUser.TimeUpdated != nil {
		gqlUser.TimeUpdated = dbUser.TimeUpdated
	}
}