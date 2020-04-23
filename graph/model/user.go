package model

import (
	"errors"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/postgres/model"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
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

func (u *User) GenerateAccessToken() (*string, error){
	if u.ID == "" {
		return nil, errors.New("user.ID is empty or invalid")
	}
	expiresAt := time.Now().Add(time.Hour * 24 * 7)
	return u.generateToken(expiresAt)
}
func (u *User) GenerateRefreshToken() (*string, error){
	if u.ID == "" {
		return nil, errors.New("user.ID is empty or invalid")
	}
	expiresAt := time.Now().Add(time.Hour * 24 * 7)
	return u.generateToken(expiresAt)
}


func (u *User) generateToken(expiresAt time.Time) (*string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: expiresAt.Unix(),
		Id:        u.ID,
		IssuedAt:  time.Now().Unix(),
		Issuer:    "innersource",
	})

	accessToken, err := token.SignedString([]byte(os.Getenv("jwt_secret")))
	if err != nil {
		return nil, err
	}
	return &accessToken, nil
}