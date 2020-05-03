package impl

import (
	"context"
	"database/sql"
	"errors"
	customErrors "github.com/cassini-Inner/inner-src-mgmt-go/custom_errors"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
	"github.com/jmoiron/sqlx"
	"strconv"
	"time"
)

var (
	ErrRemovingCurrentUserSkills = errors.New("error occurred while deleting user's existing skills")
)

type UsersRepoImpl struct {
	db *sqlx.DB
}

func (u *UsersRepoImpl) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return u.db.BeginTxx(ctx, nil)
}

func NewUsersRepo(db *sqlx.DB) *UsersRepoImpl {
	return &UsersRepoImpl{db: db}
}

func (u *UsersRepoImpl) UpdateUser(userInformation *dbmodel.User, tx *sqlx.Tx) (*dbmodel.User, error) {
	// update the users information in the database
	_, err := tx.Exec(updateUserByUserIdQuery, userInformation.Email, userInformation.Name, userInformation.Role, userInformation.Department, userInformation.Bio, userInformation.PhotoUrl, userInformation.Contact, time.Now(), userInformation.IsDeleted, userInformation.GithubUrl, userInformation.Onboarded, userInformation.GithubId, userInformation.GithubName, userInformation.Id)
	if err != nil {
		return nil, err
	}
	// fetch the updated data from the db and return
	return u.GetByIdTx(userInformation.Id, tx)
}

func (u *UsersRepoImpl) RemoveUserSkillsByUserId(userID string, tx *sqlx.Tx) error {
	_, err := tx.Exec(deleteSkillsFromUserskillsByUserIdQuery, userID)
	if err != nil {
		return ErrRemovingCurrentUserSkills
	}
	return nil
}

func (u *UsersRepoImpl) GetByIdTx(userId string, tx *sqlx.Tx) (*dbmodel.User, error) {
	_, err := strconv.Atoi(userId)
	if err != nil {
		return nil, customErrors.ErrInvalidId
	}
	var user dbmodel.User
	err = tx.QueryRowx(selectUserByIdQuery, userId).StructScan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UsersRepoImpl) GetById(userId string) (*dbmodel.User, error) {
	var user dbmodel.User
	err := u.db.QueryRowx(selectUserByIdQuery, userId).StructScan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UsersRepoImpl) GetByEmailId(emailId string) (*dbmodel.User, error) {
	var user dbmodel.User
	err := u.db.QueryRowx(selectUsersByEmailIdQuery, emailId).StructScan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UsersRepoImpl) GetByGithubId(githubId string) (*dbmodel.User, error) {
	var user dbmodel.User
	err := u.db.QueryRowx(SelectUsersByGithubIdQuery, githubId).StructScan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UsersRepoImpl) CreateNewUser(user *dbmodel.User, tx *sqlx.Tx) (*dbmodel.User, error) {
	newUserId := 0
	// we are setting up a users bio on sign-up since it's not included as part of on-boarding
	err := tx.QueryRowx(createNewUserQuery, user.Email, user.Name, user.Role, user.Department, user.Bio, user.PhotoUrl, user.Contact, user.GithubUrl, user.GithubId, user.GithubName).Scan(&newUserId)
	if err != nil {
		return nil, err
	}
	createdUser, err := u.GetByIdTx(strconv.Itoa(newUserId), tx)
	if err != nil {
		return nil, err
	}
	return createdUser, nil
}

func (u *UsersRepoImpl) CountUsersByGithubId(githubId sql.NullString, tx *sqlx.Tx) (int, error) {
	usersCount := 0
	err := tx.QueryRowx(countUsersByGithubIdQuery, githubId).Scan(&usersCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, customErrors.ErrNoEntityMatchingId
		}
		return 0, err
	}
	return usersCount, nil
}

const (
	selectUserByIdQuery        = `select * from users where users.id = $1 and users.is_deleted = false`
	selectUsersByEmailIdQuery  = `select * from users where email = $1 and users.is_deleted = false`
	SelectUsersByGithubIdQuery = `select * from users where github_id = $1 and users.is_deleted = false`
	countUsersByGithubIdQuery  = `select count(*) from users where github_id = $1`
	createNewUserQuery         = `insert into users(email, name, role, department, bio, photo_url, contact, github_url, github_id, github_name) VALUES ($1, $2, $3, $4,$5, $6, $7, $8, $9, $10) returning id`
	updateUserByUserIdQuery    = `update users set email = $1, name = $2, role = $3, department = $4, bio = $5, photo_url = $6, contact = $7, time_updated = $8, is_deleted = $9, github_url = $10, onboarded = $11, github_id = $12, github_name = $13 where id = $14`
)
