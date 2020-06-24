package impl

import (
	"context"
	"database/sql"
	"github.com/cassini-Inner/inner-src-mgmt-go/custom_errors"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/cassini-Inner/inner-src-mgmt-go/middleware"
	"github.com/cassini-Inner/inner-src-mgmt-go/repository"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
	"strings"
)

type ApplicationsService struct {
	jobsRepo          repository.JobsRepo
	applicationsRepo  repository.ApplicationsRepo
	milestonesRepo    repository.MilestonesRepo
	notificationsRepo repository.NotificationsRepo
}

func NewApplicationsService(jobsRepo repository.JobsRepo, applicationsRepo repository.ApplicationsRepo, milestonesRepo repository.MilestonesRepo, notificationRepo repository.NotificationsRepo) *ApplicationsService {
	return &ApplicationsService{
		jobsRepo:          jobsRepo,
		applicationsRepo:  applicationsRepo,
		milestonesRepo:    milestonesRepo,
		notificationsRepo: notificationRepo,
	}
}

// creates a job application of currently signed in user. This applies to a whole job
// TODO: allow people to apply to individual milestones
func (a *ApplicationsService) CreateUserJobApplication(ctx context.Context, jobId string) ([]*gqlmodel.Application, error) {
	user, err := middleware.GetCurrentUserFromContext(ctx)

	if err != nil {
		return nil, custom_errors.ErrUserNotAuthenticated
	}

	tx, err := a.applicationsRepo.BeginTx(ctx)

	job, err := a.jobsRepo.GetById(jobId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.ErrNoEntityMatchingId
		}
		return nil, err
	}

	if job.Status == "completed" {
		return nil, custom_errors.ErrJobAlreadyCompleted
	}

	if job.CreatedBy == user.Id {
		return nil, custom_errors.ErrOwnerApplyToOwnJob
	}

	milestones, err := a.milestonesRepo.GetByJobId(tx, jobId)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	totalUnassignedCount := 0
	for _, milestone := range milestones {
		// if a milestone is not assigned to anyone field is not valid as the table value is null
		if !milestone.AssignedTo.Valid {
			totalUnassignedCount++
		}
	}

	if totalUnassignedCount != len(milestones) {
		return nil, custom_errors.ErrJobAlreadyAssigned
	}

	// get all the pending or accepted applications of a user
	existingApplications, err := a.applicationsRepo.GetExistingUserApplications(tx, milestones, user.Id, dbmodel.ApplicationStatusPending, dbmodel.ApplicationStatusAccepted)

	// if some error occurred
	if err != nil && err != custom_errors.ErrNoExistingApplications {
		_ = tx.Rollback()
		return nil, err
	}

	// if no applications exist where status = pending or accepted
	if err == custom_errors.ErrNoExistingApplications {
		// get applications of user that were withdrawn or rejected
		existingApplications, err = a.applicationsRepo.GetExistingUserApplications(tx, milestones, user.Id, dbmodel.ApplicationStatusWithdrawn, dbmodel.ApplicationStatusRejected)
		if err != nil && err != custom_errors.ErrNoExistingApplications {
			_ = tx.Rollback()
			return nil, err
		}
		if err == custom_errors.ErrNoExistingApplications {
			// if the user hasn't applied to the job ever, create new application
			createdApplications, err := a.applicationsRepo.CreateApplication(ctx, tx, milestones, user.Id)
			if err != nil {
				_ = tx.Rollback()
				return nil, err
			}

			_, err = a.notificationsRepo.CreateWithTx(tx, job.CreatedBy, user.Id, dbmodel.NotificationTypeApplicationCreated, job.Id)
			if err != nil {
				err = tx.Rollback()
				if err != nil {
					return nil, err
				}
				return nil, err
			}

			err = tx.Commit()
			if err != nil {
				return nil, err
			}
			return gqlmodel.MapDBApplicationListToGql(createdApplications), nil
		}
		note := ""
		existingApplications, err := a.applicationsRepo.SetApplicationStatusForUserAndJob(ctx, tx, milestones, dbmodel.ApplicationStatusPending, &note, jobId, user.Id)

		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
		_, err = a.notificationsRepo.CreateWithTx(tx, job.CreatedBy, user.Id, dbmodel.NotificationTypeApplicationCreated, job.Id)
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				return nil, err
			}
			return nil, err
		}

		err = tx.Commit()
		if err != nil {
			return nil, err
		}
		return gqlmodel.MapDBApplicationListToGql(existingApplications), nil
	}
	_, err = a.notificationsRepo.CreateWithTx(tx, job.CreatedBy, user.Id, dbmodel.NotificationTypeApplicationCreated, job.Id)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return gqlmodel.MapDBApplicationListToGql(existingApplications), nil
}

func (a *ApplicationsService) DeleteUserJobApplication(ctx context.Context, jobId string) ([]*gqlmodel.Application, error) {
	user, err := middleware.GetCurrentUserFromContext(ctx)
	if err != nil {
		return nil, custom_errors.ErrUserNotAuthenticated
	}
	tx, err := a.applicationsRepo.BeginTx(ctx)

	if err != nil {
		return nil, err
	}

	jobMilestones, err := a.milestonesRepo.GetByJobId(tx, jobId)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// if the user is withdrawing an application where they were previously accepted into a job/milestone
	for _, milestone := range jobMilestones {
		if milestone.AssignedTo.Valid && milestone.AssignedTo.String == user.Id {
			_, err = a.milestonesRepo.SetMilestoneAssignedTo(tx, milestone.Id, nil)
			if err != nil {
				return nil, err
			}
		}
	}

	applications, err := a.applicationsRepo.SetApplicationStatusForUserAndJob(ctx, tx, jobMilestones, dbmodel.ApplicationStatusWithdrawn, nil, jobId, user.Id)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	var result []*gqlmodel.Application

	for _, application := range applications {
		var temp gqlmodel.Application
		temp.MapDbToGql(application)
		result = append(result, &temp)
	}
	return result, nil
}

func (a *ApplicationsService) UpdateJobApplicationStatus(ctx context.Context, applicantId string, jobId string, status *gqlmodel.ApplicationStatus, note *string) ([]*gqlmodel.Application, error) {
	// since this end point can only be user by job owner,
	// they can only modify job status from pending to accepted or pending
	tx, err := a.applicationsRepo.BeginTx(ctx)

	if err != nil {
		return nil, err
	}
	currentStatus, err := a.applicationsRepo.GetApplicationStatusForUserAndJob(applicantId, tx, jobId)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// check if the currently authenticate user is the owner of the job
	currentUser, err := middleware.GetCurrentUserFromContext(ctx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	currentJob, err := a.jobsRepo.GetById(jobId)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	if currentJob.CreatedBy != currentUser.Id {
		_ = tx.Rollback()
		return nil, custom_errors.ErrUserNotOwner
	}

	// owner cannot modify the status of application what was withdrawn by applicant
	// owner can only move an application from
	// - pending->accepted, pending->rejected, accepted->rejected
	if currentStatus == dbmodel.ApplicationStatusWithdrawn || currentStatus == dbmodel.ApplicationStatusRejected {
		_ = tx.Rollback()
		return nil, custom_errors.ErrApplicationWithdrawnOrRejected
	}
	// owner cannot move the application from pending or withdrawn state to any new state
	if status.String() == "PENDING" || status.String() == "WITHDRAWN" {
		_ = tx.Rollback()
		return nil, custom_errors.ErrInvalidNewApplicationState
	}

	milestones, err := a.milestonesRepo.GetByJobId(tx, jobId)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// if a new user is being assigned to a job then check if there is someone
	// already assigned to that job
	if status.String() == "ACCEPTED" {
		for _, milestone := range milestones {
			if milestone.AssignedTo.Valid {
				err := tx.Rollback()
				if err != nil {
					return nil, err
				}
				return nil, custom_errors.ErrMilestoneAlreadyAssignedOrCompleted
			}
		}
	}

	updateResult, err := a.applicationsRepo.SetApplicationStatusForUserAndJob(ctx, tx, milestones, strings.ToLower(status.String()), note, jobId, applicantId)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	var idToBeAssigned *string
	if status.String() == "ACCEPTED" {
		idToBeAssigned = &applicantId
	} else {
		idToBeAssigned = nil
	}

	updatedMilestones := make([]*dbmodel.Milestone, 0)
	for _, milestone := range milestones {
		updatedMilestone, err := a.milestonesRepo.SetMilestoneAssignedTo(tx, milestone.Id, idToBeAssigned)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
		updatedMilestones = append(updatedMilestones, updatedMilestone)
	}

	var notificationToBeSentType string

	if status.String() == "ACCEPTED" {
		notificationToBeSentType = dbmodel.NotificationTypeApplicationAccepted
	} else {
		notificationToBeSentType = dbmodel.NotificationTypeApplicationRejected
	}
	_, err = a.notificationsRepo.CreateWithTx(tx, applicantId, currentJob.CreatedBy, notificationToBeSentType, currentJob.Id)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return gqlmodel.MapDBApplicationListToGql(updateResult), nil
}

func (a *ApplicationsService) GetApplicationStatusForUserAndJob(ctx context.Context, userId string, joinId string) (string, error) {
	tx, err := a.applicationsRepo.BeginTx(ctx)
	if err != nil {
		return "", err
	}
	return a.applicationsRepo.GetApplicationStatusForUserAndJob(userId, tx, joinId)
}

func (a *ApplicationsService) GetAppliedJobs(ctx context.Context, userId string) ([]*dbmodel.Job, error) {
	return a.applicationsRepo.GetUserJobApplications(userId)
}

func (a *ApplicationsService) GetAcceptedAppliedJobs(ctx context.Context, userId string) ([]*dbmodel.Job, error) {
	return a.applicationsRepo.GetUserAcceptedJobApplications(userId)
}
