package impl

import (
	"context"
	"database/sql"
	"github.com/cassini-Inner/inner-src-mgmt-go/custom_errors"
	"github.com/cassini-Inner/inner-src-mgmt-go/middleware"
	"github.com/cassini-Inner/inner-src-mgmt-go/repository"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
)

type ReviewsService struct {
	reviewsRepo    repository.ReviewsRepo
	jobsRepo       repository.JobsRepo
	milestonesRepo repository.MilestonesRepo
}

func NewReviewsService(reviewsRepo repository.ReviewsRepo, jobsRepo repository.JobsRepo, milestonesRepo repository.MilestonesRepo) *ReviewsService {
	return &ReviewsService{reviewsRepo: reviewsRepo, jobsRepo: jobsRepo, milestonesRepo: milestonesRepo}
}

func (r ReviewsService) ReviewAssignedUser(ctx context.Context, rating int, remark *string, milestoneId string) (*dbmodel.Review, error) {
	currentRequestUser, err := middleware.GetCurrentUserFromContext(ctx)
	if err != nil {
		return nil, custom_errors.ErrUserNotAuthenticated
	}
	valid, err := r.validatedReviewData(currentRequestUser.Id, milestoneId, rating)
	if !valid {
		return nil, err
	}

	milestone, err := r.milestonesRepo.GetById(milestoneId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.ErrNoEntityMatchingId
		}
		return nil, err
	}
	if !milestone.AssignedTo.Valid {
		return nil, custom_errors.ErrNoUserAssigned
	}
	if milestone.Status != "completed" {
		return nil, custom_errors.ErrMilestoneNotCompleted
	}

	tx, err := r.reviewsRepo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	newReview := dbmodel.Review{Rating: rating, UserId: milestone.AssignedTo.String, MilestoneId: milestone.Id}
	if remark != nil {
		newReview.Remark = *remark
	}
	createdReview, err := r.reviewsRepo.Add(tx, newReview)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	err = r.reviewsRepo.CommitTx(ctx, tx)
	if err != nil {
		return nil, err
	}
	return createdReview, nil
}

func (r ReviewsService) UpdateReview(ctx context.Context, rating int, remark *string, reviewId string) (*dbmodel.Review, error) {
	currentRequestUser, err := middleware.GetCurrentUserFromContext(ctx)
	if err != nil {
		return nil, custom_errors.ErrUserNotAuthenticated
	}
	review, err := r.reviewsRepo.GetById(reviewId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.ErrInvalidId
		}
		return nil, err
	}
	milestoneAuthor, err := r.milestonesRepo.GetAuthor(review.MilestoneId)
	if err != nil {
		return nil, err
	}
	if milestoneAuthor.Id != currentRequestUser.Id {
		return nil, custom_errors.ErrUserNotOwner
	}

	if rating < 1 || rating > 5 {
		return nil, custom_errors.ErrInvalidRating
	}
	review.Rating = rating
	if remark != nil {
		review.Remark = *remark
	}
	tx, err := r.reviewsRepo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	updatedReview, err := r.reviewsRepo.Update(tx, reviewId, *review)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	err = r.reviewsRepo.CommitTx(ctx, tx)
	if err != nil {
		return nil, err
	}
	return updatedReview, nil
}

func (r ReviewsService) validatedReviewData(userId, milestoneId string, rating int) (bool, error) {
	milestoneAuthor, err := r.milestonesRepo.GetAuthor(milestoneId)
	if err != nil {
		return false, err
	}
	if userId != milestoneAuthor.Id {
		return false, custom_errors.ErrUserNotOwner
	}
	if milestoneId == "" {
		return false, custom_errors.ErrInvalidId
	}
	if rating < 1 || rating > 5 {
		return false, custom_errors.ErrInvalidRating
	}
	return true, nil
}

func (r ReviewsService) GetForUserId(ctx context.Context, userId string) ([]*dbmodel.Review, error) {
	return r.reviewsRepo.GetForUserId(userId)
}
