package model

import(
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/postgres/model"
)


type Application struct {
	ID          string            `json:"id"`
	ApplicantID string            `json:"applicant"`
	Status      ApplicationStatus `json:"status"`
	Note        *string           `json:"note"`
	CreatedOn   string            `json:"createdOn"`
}

func (a *Application) mapDbToGql(dbApplication dbmodel.Application) {
	a.ID = dbApplication.Id
	a.ApplicantID = dbApplication.ApplicantId
	a.Status = ApplicationStatus(dbApplication.Status)
	a.Note = &dbApplication.Note
	a.CreatedOn = dbApplication.TimeCreated
}