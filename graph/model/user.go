package model

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
