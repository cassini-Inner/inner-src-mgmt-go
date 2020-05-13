package model

import (
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
)

type Comment struct {
	ID          string `json:"id"`
	TimeCreated string `json:"timeCreated"`
	TimeUpdated string `json:"timeUpdated"`
	Content     string `json:"content"`
	IsDeleted   bool   `json:"isDeleted"`
	CreatedBy   string `json:"createdBy"`
}

func (c *Comment) MapDbToGql(dbComment dbmodel.Discussion) {
	c.ID = dbComment.Id
	c.TimeCreated = dbComment.TimeCreated
	c.TimeCreated = dbComment.TimeUpdated
	c.Content = dbComment.Content
	c.CreatedBy = dbComment.CreatedBy
	c.IsDeleted = dbComment.IsDeleted
}
