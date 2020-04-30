package impl

import (
	"database/sql"
	"errors"
	customErrors "github.com/cassini-Inner/inner-src-mgmt-go/custom_errors"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/postgres/model"
	"github.com/jmoiron/sqlx"
	"log"
	"strconv"
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

func (u *UsersRepoImpl) UpdateUser(currentUserInfo *dbmodel.User, input *gqlmodel.UpdateUserInput, tx *sqlx.Tx) (*dbmodel.User, error) {

	updatedUserInformation := *currentUserInfo
	// setup a new transaction
	// check which fields need to be updated
	if input.Contact != nil {
		updatedUserInformation.Contact = dbmodel.ToNullString(input.Contact)
	}
	if input.Bio != nil {
		updatedUserInformation.Bio = dbmodel.ToNullString(input.Bio)
	}
	if input.Department != nil {
		updatedUserInformation.Department = dbmodel.ToNullString(input.Department)
	}
	if input.Role != nil {
		updatedUserInformation.Role = dbmodel.ToNullString(input.Role)
	}
	if input.Name != nil {
		updatedUserInformation.Name = dbmodel.ToNullString(input.Name)
	}
	if input.Email != nil {
		updatedUserInformation.Email = dbmodel.ToNullString(input.Email)
	}

	// update the users information in the database
	_, err := tx.Exec(updateUserByUserIdQuery, updatedUserInformation.Email, updatedUserInformation.Name, updatedUserInformation.Role, updatedUserInformation.Department, updatedUserInformation.Bio, updatedUserInformation.Contact, updatedUserInformation.Id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// delete current entries from userskills table
	if input.Skills != nil {

	}

	return &updatedUserInformation, nil
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
	updateUserByUserIdQuery    = `update users set email = $1, name = $2, role = $3, department = $4, bio = $5, contact = $6, onboarded=true where id = $7`
)
