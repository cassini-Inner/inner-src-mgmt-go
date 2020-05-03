package service

import (
	"context"
	"github.com/cassini-Inner/inner-src-mgmt-go/custom_errors"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/cassini-Inner/inner-src-mgmt-go/middleware"
	"github.com/cassini-Inner/inner-src-mgmt-go/repository"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
	"strconv"
)

type UserProfileService struct {
	userRepo   repository.UsersRepo
	skillsRepo repository.SkillsRepo
}

func NewUserProfileService( userRepo repository.UsersRepo, skillsRepo repository.SkillsRepo) *UserProfileService {
	return &UserProfileService{ userRepo: userRepo, skillsRepo: skillsRepo}
}

func (s UserProfileService) UpdateProfile(ctx context.Context, userDetails *gqlmodel.UpdateUserInput) (*gqlmodel.User, error) {
	currentRequestUser, err := middleware.GetCurrentUserFromContext(ctx)
	if err != nil {
		return nil, custom_errors.ErrUserNotAuthenticated
	}

	tx, err := s.skillsRepo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	if userDetails.Contact != nil {
		currentRequestUser.Contact = dbmodel.ToNullString(userDetails.Contact)
	}
	if userDetails.Bio != nil {
		currentRequestUser.Bio = dbmodel.ToNullString(userDetails.Bio)
	}
	if userDetails.Department != nil {
		currentRequestUser.Department = dbmodel.ToNullString(userDetails.Department)
	}
	if userDetails.Role != nil {
		currentRequestUser.Role = dbmodel.ToNullString(userDetails.Role)
	}
	if userDetails.Email != nil {
		currentRequestUser.Email = dbmodel.ToNullString(userDetails.Email)
	}
	if userDetails.Name != nil {
		currentRequestUser.Name = dbmodel.ToNullString(userDetails.Name)
	}
	user, err := s.userRepo.UpdateUser(currentRequestUser, tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
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
		return nil, custom_errors.ErrInvalidId
	}
	user, err := s.userRepo.GetById(userId)
	if err != nil {
		return nil, err
	}
	var gqlUser gqlmodel.User
	gqlUser.MapDbToGql(*user)
	return &gqlUser, nil
}
