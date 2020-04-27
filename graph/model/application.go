package model

import (
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/postgres/model"
	"strings"
)

type Application struct {
	ID          string            `json:"id"`
	ApplicantID string            `json:"applicant"`
	MilestoneID string            `json:"applicant"`
	Status      ApplicationStatus `json:"status"`
	Note        *string           `json:"note"`
	CreatedOn   string            `json:"createdOn"`
}

func (a *Application) MapDbToGql(dbApplication *dbmodel.Application) {
	a.ID = dbApplication.Id
	a.ApplicantID = dbApplication.ApplicantId
	a.Status = ApplicationStatus(strings.ToUpper(dbApplication.Status))
	if dbApplication.Note.Valid {
		a.Note = &dbApplication.Note.String
	}
	a.MilestoneID = dbApplication.MilestoneId
	a.CreatedOn = dbApplication.TimeCreated
}

func MapDBApplicationListToGql(dbApplicationList []*dbmodel.Application) []*Application {
	var result []*Application
	for _, dbApplication := range dbApplicationList {
		var a Application
		a.ID = dbApplication.Id
		a.ApplicantID = dbApplication.ApplicantId
		a.Status = ApplicationStatus(strings.ToUpper(dbApplication.Status))
		if dbApplication.Note.Valid {
			a.Note = &dbApplication.Note.String
		}
		a.MilestoneID = dbApplication.MilestoneId
		a.CreatedOn = dbApplication.TimeCreated

		result = append(result, &a)
	}

	return result
}
