package impl

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	customErrors "github.com/cassini-Inner/inner-src-mgmt-go/custom_errors"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
	"github.com/jmoiron/sqlx"
)

var (
	ErrRemovingCurrentUserSkills = errors.New("error occurred while deleting user's existing skills")
)

type UsersRepoImpl struct {
	db *sqlx.DB
}

func NewUsersRepo(db *sqlx.DB) *UsersRepoImpl {
	return &UsersRepoImpl{db: db}
}

func (u *UsersRepoImpl) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return u.db.BeginTxx(ctx, nil)
}

func (u *UsersRepoImpl) CommitTx(ctx context.Context, tx *sqlx.Tx) (err error) {
	err = tx.Commit()
	if err != nil {
		err = tx.Rollback()
	}
	return err
}

func (u *UsersRepoImpl) UpdateUser(tx *sqlx.Tx, updatedUserInformation *dbmodel.User) (*dbmodel.User, error) {
	// update the users information in the database
	_, err := tx.Exec(updateUserByUserIdQuery, updatedUserInformation.Email, updatedUserInformation.Name, updatedUserInformation.Role, updatedUserInformation.Department, updatedUserInformation.Bio, updatedUserInformation.PhotoUrl, updatedUserInformation.Contact, time.Now(), updatedUserInformation.IsDeleted, updatedUserInformation.GithubUrl, updatedUserInformation.Onboarded, updatedUserInformation.GithubId, updatedUserInformation.GithubName, updatedUserInformation.Id)
	if err != nil {
		return nil, err
	}

	// fetch the updated data from the db and return
	return u.GetByIdTx(tx, updatedUserInformation.Id)
}

func (u *UsersRepoImpl) RemoveUserSkillsByUserId(tx *sqlx.Tx, userID string) error {
	_, err := tx.Exec(deleteSkillsFromUserskillsByUserIdQuery, userID)
	if err != nil {
		return ErrRemovingCurrentUserSkills
	}
	return nil
}

func (u *UsersRepoImpl) GetByIdTx(tx *sqlx.Tx, userId string) (*dbmodel.User, error) {
	_, err := strconv.Atoi(userId)
	if err != nil {
		return nil, customErrors.ErrInvalidId
	}
	user := &dbmodel.User{}
	err = tx.QueryRowx(selectUserByIdQuery, userId).StructScan(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UsersRepoImpl) GetById(userId string) (*dbmodel.User, error) {
	user := &dbmodel.User{}
	err := u.db.QueryRowx(selectUserByIdQuery, userId).StructScan(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UsersRepoImpl) GetByEmailId(emailId string) (*dbmodel.User, error) {
	user := &dbmodel.User{}
	err := u.db.QueryRowx(selectUsersByEmailIdQuery, emailId).StructScan(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UsersRepoImpl) GetByGithubId(githubId string) (*dbmodel.User, error) {
	user := &dbmodel.User{}
	err := u.db.QueryRowx(SelectUsersByGithubIdQuery, githubId).StructScan(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UsersRepoImpl) CreateNewUser(tx *sqlx.Tx, user *dbmodel.User) (*dbmodel.User, error) {
	newUserId := 0
	// we are setting up a users bio on sign-up since it's not included as part of on-boarding
	err := tx.QueryRowx(createNewUserQuery, user.Email, user.Name, user.Role, user.Department, user.Bio, user.PhotoUrl, user.Contact, user.GithubUrl, user.GithubId, user.GithubName).Scan(&newUserId)
	if err != nil {
		return nil, err
	}
	createdUser, err := u.GetByIdTx(tx, strconv.Itoa(newUserId))
	if err != nil {
		return nil, err
	}
	return createdUser, nil
}

func (u *UsersRepoImpl) CountUsersByGithubId(tx *sqlx.Tx, githubId sql.NullString) (int, error) {
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

func (u *UsersRepoImpl) GetByName(userName string, limit *int) ([]dbmodel.User, error) {
	rows, err := u.db.Queryx(selectUserNameQuery, userName, limit)
	if err != nil {
		return nil, err
	}

	var users []dbmodel.User
	for rows != nil && rows.Next() {
		tempUser := &dbmodel.User{}
		rows.StructScan(tempUser)
		users = append(users, *tempUser)
	}

	return users, nil
}

const (
	selectUserByIdQuery        = `select * from users where users.id = $1 and users.is_deleted = false`
	selectUserNameQuery        = `select * from users where users.name ~* $1 and users.is_deleted = false order by users.name limit $2`
	selectUsersByEmailIdQuery  = `select * from users where email = $1 and users.is_deleted = false`
	SelectUsersByGithubIdQuery = `select * from users where github_id = $1 and users.is_deleted = false`
	countUsersByGithubIdQuery  = `select count(*) from users where github_id = $1`
	createNewUserQuery         = `insert into users(email, name, role, department, bio, photo_url, contact, github_url, github_id, github_name) VALUES ($1, $2, $3, $4,$5, $6, $7, $8, $9, $10) returning id`
	updateUserByUserIdQuery    = `update users set email = $1, name = $2, role = $3, department = $4, bio = $5, photo_url = $6, contact = $7, time_updated = $8, is_deleted = $9, github_url = $10, onboarded = $11, github_id = $12, github_name = $13 where id = $14`
)
