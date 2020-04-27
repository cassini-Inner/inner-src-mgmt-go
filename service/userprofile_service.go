package service

import (
	"context"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/cassini-Inner/inner-src-mgmt-go/middleware"
	"github.com/cassini-Inner/inner-src-mgmt-go/postgres"
	"github.com/jmoiron/sqlx"
	"strconv"
)

type UserProfileService struct {
	db         *sqlx.DB
	userRepo   *postgres.UsersRepo
	skillsRepo *postgres.SkillsRepo
}

func NewUserProfileService(db *sqlx.DB, userRepo *postgres.UsersRepo, skillsRepo *postgres.SkillsRepo) *UserProfileService {
	return &UserProfileService{db: db, userRepo: userRepo, skillsRepo: skillsRepo}
}

func (s UserProfileService) UpdateProfile(ctx context.Context, userDetails *gqlmodel.UpdateUserInput) (*gqlmodel.User, error) {
	currentRequestUser, err := middleware.GetCurrentUserFromContext(ctx)
	if err != nil {
		return nil, ErrUserNotAuthenticated
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}


	// if any of these fields are being updated
	if userDetails.Contact != nil || userDetails.Bio != nil || userDetails.Department != nil || userDetails.Name != nil || userDetails.Email != nil {
		_, err := s.userRepo.UpdateUser(currentRequestUser, userDetails, tx)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	}

	// if user skills are to be created
	if userDetails.Skills != nil {
		err := s.userRepo.RemoveUserSkillsByUserId(currentRequestUser.Id, tx)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}

		if len(userDetails.Skills) != 0 {
			//build the slice of skills for user
			skillList := make([]string, len(userDetails.Skills))
			for i, skill := range userDetails.Skills {
				skillList[i] = *skill
			}

			// create new skills for the users
			newSkills, err := s.skillsRepo.FindOrCreateSkills(ctx, tx, skillList, currentRequestUser.Id)
			if err != nil {
				_ = tx.Rollback()
				return nil, err
			}
			err = s.skillsRepo.AddSkillsToUserSkills(newSkills, tx, currentRequestUser.Id)
			if err != nil {
				_ = tx.Rollback()
				return nil, err
			}
		}
	}

	user, err := s.userRepo.GetByIdTx(currentRequestUser.Id, tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	var gqlUpdatedUser gqlmodel.User
	gqlUpdatedUser.MapDbToGql(*user)

	return &gqlUpdatedUser, nil
}

func (s *UserProfileService) GetById(ctx context.Context, userId string) (*gqlmodel.User, error) {
	if _, err := strconv.Atoi(userId); err != nil || userId == "" {
		return nil, ErrInvalidId
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepo.GetByIdTx(userId, tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	var gqlUser gqlmodel.User
	gqlUser.MapDbToGql(*user)
	return &gqlUser, nil
}
