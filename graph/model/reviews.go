package model

import dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"

type JobReview struct {
	Job             *Job               `json:"job"`
	MilestoneReview []*MilestoneReview `json:"milestoneReview"`
}

type MilestoneReview struct {
	Review    *Review    `json:"review"`
	Milestone *Milestone `json:"milestone"`
}

type Review struct {
	Id          string  `json:"id"`
	Rating      int     `json:"rating"`
	Remark      *string `json:"remark"`
	CreatedFor  string  `json:"created_for"`
	TimeCreated string  `json:"timeCreated"`
	TimeUpdated string  `json:"timeUpdates"`
}

type ReviewInput struct {
	Rating int     `json:"rating"`
	Remark *string `json:"remark"`
}

func (r *Review) MapDbToGql(review dbmodel.Review) {
	r.Id = review.Id
	r.Remark = &review.Remark
	r.Rating = review.Rating
	r.TimeCreated = review.TimeCreated
	r.TimeUpdated = review.TimeUpdated
	r.CreatedFor = review.UserId
}
