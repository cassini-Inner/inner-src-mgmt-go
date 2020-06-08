package service

import (
	"context"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
)

type ReviewsService interface {
	// adds a review to the user that's currently assigned to a milestone
	ReviewAssignedUser(ctx context.Context, rating int, remark *string, milestoneId string) (*dbmodel.Review, error)
	// update a review based on it's ID
	UpdateReview(ctx context.Context, rating int, remark *string, reviewId string) (*dbmodel.Review, error)
}
