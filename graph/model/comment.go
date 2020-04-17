package model

import(
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/postgres/models"
)

type Comment struct {
	ID          string `json:"id"`
	TimeCreated string `json:"timeCreated"`
	TimeUpdated string `json:"timeUpdated"`
	Content     string `json:"content"`
	IsDeleted   bool   `json:"isDeleted"`
	CreatedBy   string `json:"createdBy"`
}

func (gqlCommnet *Commnet) mapDbToGql(dbCommnet dbmodel.discussion) {

	if dbCommnet.Id != nil {
		gqlCommnet.ID = dbCommnet.Id
	}

	if dbCommnet.TimeCreated != nil {
		gqlCommnet.TimeCreated = dbCommnet.TimeCreated
	}

	if dbCommnet.TimeUpdated != nil {
		gqlCommnet.TimeUpdated = dbCommnet.TimeUpdated
	}

	if dbCommnet.Content != nil {
		gqlCommnet.Content = dbCommnet.Content
	}

	if dbCommnet.CreatedBy != nil {
		gqlCommnet.CreatedBy = dbCommnet.CreatedBy
	}

	if dbCommnet.IsDeleted != nil {
		gqlCommnet.IsDeleted = dbCommnet.IsDeleted
	}
}