package service

import (
	"context"
	"database/sql"
	"errors"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/cassini-Inner/inner-src-mgmt-go/middleware"
	"github.com/cassini-Inner/inner-src-mgmt-go/postgres"
	"github.com/jmoiron/sqlx"
	"log"
)

var (
	ErrUserNotOwner                   = errors.New("current user is not owner of this entity, and hence cannot modify it")
	ErrNoEntityMatchingId             = errors.New("no entity found that matches given id")
	ErrOwnerApplyToOwnJob             = errors.New("owner cannot apply to their job")
	ErrApplicationWithdrawnOrRejected = errors.New("owner cannot modify applications with withdrawn status")
	ErrInvalidNewApplicationState     = errors.New("owner cannot move application status to withdrawn or pending")
	ErrJobAlreadyCompleted            = errors.New("job is already completed")
	ErrEntityDeleted                  = errors.New("entity was deleted")
	ErrUserNotAuthenticated           = errors.New("unauthorized request")
)

type JobsService struct {
	db               *sqlx.DB
	jobsRepo         *postgres.JobsRepo
	skillsRepo       *postgres.SkillsRepo
	discussionsRepo  *postgres.DiscussionsRepo
	applicationsRepo *postgres.ApplicationsRepo
}

func NewJobsService(db *sqlx.DB, jobsRepo *postgres.JobsRepo, skillsRepo *postgres.SkillsRepo, discussionsRepo *postgres.DiscussionsRepo, applicationsRepo *postgres.ApplicationsRepo) *JobsService {
	return &JobsService{db: db, jobsRepo: jobsRepo, skillsRepo: skillsRepo, discussionsRepo: discussionsRepo, applicationsRepo:applicationsRepo}
}

func (j *JobsService) CreateJob(ctx context.Context, job *gqlmodel.CreateJobInput) (result *gqlmodel.Job, err error) {
	// validate the input
	if len(job.Desc) < 5 {
		return nil, errors.New("description not long enough")
	}
	if len(job.Title) < 5 {
		return nil, errors.New("title not long enough")
	}
	if len(job.Difficulty) == 5 {
		return nil, errors.New("diff not long enough")
	}
	if len(job.Milestones) == 0 {
		return nil, errors.New("just must have at least one milestone")
	}

	for _, milestone := range job.Milestones {
		if len(milestone.Skills) == 0 {
			return nil, errors.New("milestone must have at least one skill")
		}
	}

	user, err := middleware.GetCurrentUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := j.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	newJob, err := j.jobsRepo.CreateJob(ctx, tx, job, user)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	newMilestones, err := j.jobsRepo.CreateMilestones(ctx, tx, newJob.Id, job.Milestones)
	if err != nil {
		return nil, err
	}

	var newSkillsList []string
	for _, milestone := range job.Milestones {
		for _, s := range milestone.Skills {
			val := *s
			newSkillsList = append(newSkillsList, val)
		}
	}

	newSkills, err := j.skillsRepo.FindOrCreateSkills(ctx, tx, newSkillsList, user.Id)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	err = j.skillsRepo.MapSkillsToMilestones(ctx, tx, newSkills, job, newMilestones)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	result = &gqlmodel.Job{}
	result.MapDbToGql(*newJob)
	return result, nil
}

func (j *JobsService) AddDiscussionToJob(ctx context.Context, comment, jobId string) (*gqlmodel.Comment, error) {
	tx, err := j.db.BeginTxx(ctx, nil)
	user, err := middleware.GetCurrentUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	newComment, err := j.discussionsRepo.CreateComment(jobId, comment, user.Id, tx, ctx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	var gqlComment gqlmodel.Comment
	gqlComment.MapDbToGql(*newComment)
	return &gqlComment, nil
}

func (j *JobsService) UpdateJobDiscussion(ctx context.Context, commentId, comment string) (*gqlmodel.Comment, error) {
	tx, err := j.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, nil
	}
	user, err := middleware.GetCurrentUserFromContext(ctx)
	if err != nil {
		return nil, ErrUserNotAuthenticated
	}

	existingDiscussion, err := j.discussionsRepo.GetById(commentId, tx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoEntityMatchingId
		}
		return nil, err
	}
	if existingDiscussion.CreatedBy != user.Id {
		return nil, ErrUserNotOwner
	}

	updatedDiscussion, err := j.discussionsRepo.UpdateComment(commentId, comment, tx, ctx)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	var gqlUpdatedDiscussion gqlmodel.Comment
	gqlUpdatedDiscussion.MapDbToGql(*updatedDiscussion)
	return &gqlUpdatedDiscussion, nil
}

func (j *JobsService) DeleteJobDiscussion(ctx context.Context, commentId string) (*gqlmodel.Comment, error) {
	if commentId == "" {
		return nil, ErrNoEntityMatchingId
	}

	user, err := middleware.GetCurrentUserFromContext(ctx)
	if err != nil {
		return nil, ErrUserNotAuthenticated
	}

	tx, err := j.db.BeginTxx(ctx, nil)

	existingDiscussion, err := j.discussionsRepo.GetById(commentId, tx)
	if err != nil {
		_ = tx.Rollback()
		if err == sql.ErrNoRows {
			return nil, ErrNoEntityMatchingId
		}
		return nil, err
	}
	if existingDiscussion.CreatedBy != user.Id {
		_ = tx.Rollback()
		return nil, ErrUserNotOwner
	}
	discussion, err := j.discussionsRepo.DeleteComment(commentId, tx, ctx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	var gqlComment gqlmodel.Comment
	gqlComment.MapDbToGql(*discussion)
	return &gqlComment, nil
}

func (j *JobsService) GetById(ctx context.Context, jobId string) (*gqlmodel.Job, error) {
	tx, err := j.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	job, err := j.jobsRepo.GetById(jobId, tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	var gqlJob gqlmodel.Job
	gqlJob.MapDbToGql(*job)
	return &gqlJob, nil
}

func (j *JobsService) ToggleJobCompleted(ctx context.Context, jobID string) (*gqlmodel.Job, error) {
	user, err := middleware.GetCurrentUserFromContext(ctx)
	if err != nil {
		return nil, ErrUserNotAuthenticated
	}

	tx, err := j.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	// check if the job exists in the repo
	job, err := j.jobsRepo.GetByIdTx(jobID, tx)
	if err != nil {
		return nil, ErrNoEntityMatchingId
	}
	if job.IsDeleted {
		return nil, ErrEntityDeleted
	}

	// check if the job is being modified by the person who created it
	if job.CreatedBy != user.Id {
		return nil, ErrUserNotOwner
	}

	// get all the milestones for a job
	milestones, err := j.jobsRepo.GetMilestonesByJobId(jobID, tx)
	if err != nil {
		_ = tx.Rollback()
		log.Println(err)
		return nil, err
	}
	// mark all the milestones as completed
	milestoneIds := make([]string, len(milestones))
	for i, milestone := range milestones {
		milestoneIds[i] = milestone.Id
	}
	if job.Status != "completed" {
		_, err := j.jobsRepo.MarkJobCompleted(ctx, tx, jobID)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
		err = j.jobsRepo.MarkMilestonesCompleted(tx, ctx, milestoneIds...)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	} else {
		_, err := j.jobsRepo.ForceAutoUpdateJobStatus(ctx, tx, jobID)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
		err = j.jobsRepo.ForceAutoUpdateMilestoneStatusByJobID(ctx, tx, jobID)
		if err != nil {
			return nil, err
		}
	}
	updatedJob, err := j.jobsRepo.GetByIdTx(jobID, tx)
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	var gqlJob gqlmodel.Job
	gqlJob.MapDbToGql(*updatedJob)

	return &gqlJob, nil
}

func (j *JobsService) DeleteJob(ctx context.Context, jobID string) (*gqlmodel.Job, error) {
	user, err := middleware.GetCurrentUserFromContext(ctx)
	if err != nil {
		return nil, ErrUserNotAuthenticated
	}

	tx, err := j.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	job, err := j.jobsRepo.GetByIdTx(jobID, tx)
	if err != nil {
		return nil, ErrNoEntityMatchingId
	}

	if user.Id != job.CreatedBy {
		return nil, ErrUserNotOwner
	}

	// delete the job from db
	deletedJob, err := j.jobsRepo.DeleteJob(tx, jobID)
	if err != nil {
		return nil, err
	}

	// delete job discussions from db
	err = j.discussionsRepo.DeleteAllCommentsForJob(tx, jobID)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// delete all the applications for the job
	err = j.applicationsRepo.DeleteAllJobApplications(tx, jobID)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	err  = j.jobsRepo.DeleteMilestonesByJobId(tx, jobID)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	 err = tx.Commit()
	if err != nil {
		return nil, err
	}
	var gqlJob gqlmodel.Job
	 gqlJob.MapDbToGql(*deletedJob)
	 return &gqlJob, nil
}

// TODO: we're mixing tx and non tx queries here. refactor to only use tx queries
func (j *JobsService) ToggleMilestoneCompleted(ctx context.Context, milestoneID string) (*gqlmodel.Milestone, error) {
	user, err := middleware.GetCurrentUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	milestoneData, err := j.jobsRepo.GetMilestoneById(milestoneID)
	if err != nil {
		return nil, err
	}
	milestoneAuthor, err := j.jobsRepo.GetAuthorFromMilestoneId(milestoneID)
	if err != nil {
		return nil, err
	}

	if milestoneAuthor.Id != user.Id {
		return nil, ErrUserNotOwner
	}

	tx, err := j.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	// if the milestone status was already completed we will need to modify on the
	if milestoneData.Status != "completed" {
		err = j.jobsRepo.MarkMilestonesCompleted(tx, ctx, milestoneID)
		if err != nil {
			return nil, err
		}

		jobMilestones, err := j.jobsRepo.GetMilestonesByJobId(milestoneData.JobId, tx)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
		completedMilestonesCount := 0
		for _, milestone := range jobMilestones {
			if milestone.Status == "completed" {
				completedMilestonesCount++
			}
		}
		if completedMilestonesCount == len(jobMilestones) {
			_, err := j.jobsRepo.MarkJobCompleted(ctx, tx, milestoneData.JobId)
			if err != nil {
				_ = tx.Rollback()
				return nil, err
			}
		}
	} else {
		err = j.jobsRepo.ForceAutoUpdateMilestoneStatusByMilestoneId(ctx, tx, milestoneID)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
		_, err := j.jobsRepo.ForceAutoUpdateJobStatus(ctx, tx, milestoneData.JobId)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	updatedMilestone, err := j.jobsRepo.GetMilestoneById(milestoneID)
	if err != nil {
		return nil, err
	}

	var gqlMilestone gqlmodel.Milestone
	gqlMilestone.MapDbToGql(*updatedMilestone)
	return &gqlMilestone, nil
}
