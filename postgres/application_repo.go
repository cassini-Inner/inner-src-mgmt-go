package postgres

import (
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/jinzhu/gorm"
)

type ApplicationsRepo struct {
	db *gorm.DB
}

func NewApplicationsRepo(db *gorm.DB) *ApplicationsRepo {
	return &ApplicationsRepo{db: db}
}

func (a *ApplicationsRepo) CreateApplication(jobId string, userId string) (*model.Application, error) {
	panic("Not implemented")
}
func (a *ApplicationsRepo) UpdateApplication(applicantId, jobId string, newStatus model.ApplicationStatus) (*model.Application, error) {
	panic("Not implemented")
}

func (a *ApplicationsRepo) DeleteApplication(jobId string, userId string) (*model.Application, error) {
	panic("Not implemented")
}

func (a *ApplicationsRepo) GetByJobId(jobId string) (*model.Applications, error) {
	panic("not implemented")
}

func (a *ApplicationsRepo) GetUserJobApplications(userId string) ([]*model.Job, error) {
	panic("not implemented")
}
