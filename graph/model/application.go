package model

type Application struct {
	ID          string            `json:"id"`
	ApplicantID string            `json:"applicant"`
	Status      ApplicationStatus `json:"status"`
	Note        *string           `json:"note"`
	CreatedOn   string            `json:"createdOn"`
}
