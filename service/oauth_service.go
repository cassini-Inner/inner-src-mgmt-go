package service

import dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"

type OAuthService interface {
	Authenticate(code string) (string, error)
	GetUserInfo() (*dbmodel.User, error)
}
