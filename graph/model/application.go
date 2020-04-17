package model

import(
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/postgres/models"
)


type Application struct {
	ID          string            `json:"id"`
	ApplicantID string            `json:"applicant"`
	Status      ApplicationStatus `json:"status"`
	Note        *string           `json:"note"`
	CreatedOn   string            `json:"createdOn"`
}

func (gqlApplication *Application) mapDbToGql(dbApplication dbmodel.application) {

	if dbApplication.Id != nil{
		gqlApplication.ID = dbApplication.Id
	}

	if dbApplication.ApplicantID != nil{
		gqlApplication.ApplicantID = dbApplication.ApplicantID
	}

	if dbApplication.Status != nil{
		gqlApplication.Status = dbApplication.Status
	}

	if dbApplication.Notes != nil{
		gqlApplication.Note = dbApplication.Notes
	}

	if dbApplication.TimeCreated != nil{
		gqlApplication.CreatedOn = dbApplication.TimeCreated
	}
}