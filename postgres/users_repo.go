package postgres

import (
	"encoding/json"
	"errors"
	"fmt"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/postgres/model"
	"github.com/jmoiron/sqlx"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	ErrRemovingCurrentUserSkills = errors.New("error occurred while deleting user's existing skills")
)

// TODO: Implement
type UsersRepo struct {
	db *sqlx.DB
}

func NewUsersRepo(db *sqlx.DB) *UsersRepo {
	return &UsersRepo{db: db}
}

//TODO: Deprecate this function
func (u *UsersRepo) CreateUser(input *gqlmodel.CreateUserInput) (*dbmodel.User, error) {
	panic("not implemented")
}

//TODO: Deprecate
func (u *UsersRepo) UpdateUser(currentUserInfo *dbmodel.User, input *gqlmodel.UpdateUserInput) (*dbmodel.User, error) {

	updatedUserInformation := *currentUserInfo
	// setup a new transaction
	tx, err := u.db.Beginx()
	if err != nil {
		return nil, err
	}

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
	_, err = tx.Exec(updateUserByUserIdQuery, updatedUserInformation.Email, updatedUserInformation.Name, updatedUserInformation.Role, updatedUserInformation.Department, updatedUserInformation.Bio, updatedUserInformation.Contact, updatedUserInformation.Id)
	if err != nil {
		log.Println(err)
		_ = tx.Rollback()
		return nil, err
	}

	// delete current entries from userskills table
	if input.Skills != nil {
		_, err = tx.Exec(deleteSkillsFromUserskillsByUserIdQuery, currentUserInfo.Id)
		if err != nil {
			_ = tx.Rollback()
			log.Println(err)
			return nil, ErrRemovingCurrentUserSkills
		}
		// create new skills for the users
		var inputSkills []string
		for _, skill := range input.Skills {
			inputSkills = append(inputSkills, *skill)
		}
		newSkills, err := findOrCreateSkills(inputSkills, currentUserInfo.Id, tx)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}

		err = addSkillsToUserSkills(newSkills, tx, currentUserInfo.Id)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Println("failed to commit update Users transaction")
		return nil, err
	}

	return &updatedUserInformation, nil
}

func (u *UsersRepo) GetById(userId string) (*dbmodel.User, error) {
	var user dbmodel.User
	err := u.db.QueryRowx(selectUserByIdQuery, userId).StructScan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UsersRepo) GetByEmailId(emailId string) (*dbmodel.User, error) {
	var user dbmodel.User
	err := u.db.QueryRowx(selectUsersByEmailIdQuery, emailId).StructScan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UsersRepo) GetByGithubId(githubId string) (*dbmodel.User, error) {
	var user dbmodel.User
	err := u.db.QueryRowx(selectUsersByGithubIdQuery, githubId).StructScan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UsersRepo) AuthenticateAndGetUser(githubCode string) (*dbmodel.User, error) {
	accessToken, err := u.getAccessTokenFromCode(githubCode)
	if err != nil {
		return nil, err
	}

	fetchedUser, err := u.getUserInformationFromToken(accessToken)
	if err != nil {
		return nil, err
	}
	// check if the user is signing up for the first time
	usersCount := 0
	err = u.db.QueryRowx(countUsersByGithubIdQuery, fetchedUser.GithubId).Scan(&usersCount)
	if err != nil {
		return nil, err
	}

	switch usersCount {
	case 0:
		return u.createNewUser(fetchedUser)
	case 1:
		user, err := u.GetByGithubId(fetchedUser.GithubId.String)
		if err != nil {
			return nil, err
		}
		return user, nil
	default:
		return nil, errors.New("multiple users by the same name exist in database")
	}

}

func (u *UsersRepo) createNewUser(user *dbmodel.User) (*dbmodel.User, error) {
	newUserId := 0
	// we are setting up a users bio on sign-up since it's not included as part of on-boarding
	err := u.db.QueryRowx(createNewUserQuery, user.Email, user.Name, user.Role, user.Department, user.Bio, user.PhotoUrl, user.Contact, user.GithubUrl, user.GithubId, user.GithubName).Scan(&newUserId)
	if err != nil {
		return nil, err
	}
	createdUser, err := u.GetById(strconv.Itoa(newUserId))
	if err != nil {
		return nil, err
	}
	return createdUser, nil
}

func (u *UsersRepo) getAccessTokenFromCode(githubCode string) (string, error) {
	urlStr := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%v&client_secret=%v&code=%v", os.Getenv("client_id"),
		os.Getenv("client_secret"),
		githubCode,
	)

	client := http.Client{}
	request, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(""))
	request.Header.Add("Accept", "application/json")
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	data, err := u.parseJsonFromResponse(response.Body)
	if err != nil {
		return "", err
	}
	//TODO: Error message is ambiguous, use the error message provided by github
	if response.StatusCode != 200 {
		return "", errors.New("could not authenticate with github")
	}
	accessToken, ok := data["access_token"].(string)
	fmt.Println(accessToken)
	if !ok {
		return "", errors.New("could not get access_token from github auth response, token expired or invalid")
	}

	return accessToken, nil
}

func (u *UsersRepo) getUserInformationFromToken(accessToken string) (*dbmodel.User, error) {
	client := http.Client{}
	request, _ := http.NewRequest(http.MethodGet, "https://api.github.com/user", strings.NewReader(""))
	request.Header.Set("Authorization", fmt.Sprintf(
		"token %v", accessToken))

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	result, err := u.parseJsonFromResponse(response.Body)
	if err != nil {
		return nil, err
	}

	user := &dbmodel.User{
		Email:      dbmodel.ToNullString(result["email"]),
		Name:       dbmodel.ToNullString(result["name"]),
		Bio:        dbmodel.ToNullString(result["bio"]),
		PhotoUrl:   dbmodel.ToNullString(result["avatar_url"]),
		GithubUrl:  dbmodel.ToNullString(result["html_url"]),
		GithubId:   dbmodel.ToNullString(result["id"]),
		GithubName: dbmodel.ToNullString(result["login"]),
	}
	return user, nil
}

func (u *UsersRepo) parseJsonFromResponse(responseBody io.Reader) (map[string]interface{}, error) {
	var result map[string]interface{}
	body, _ := ioutil.ReadAll(responseBody)
	err := json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

const (
	selectUserByIdQuery        = `select * from users where users.id = $1 and users.is_deleted = false`
	selectUsersByEmailIdQuery  = `select * from users where email = $1 and users.is_deleted = false`
	selectUsersByGithubIdQuery = `select * from users where github_id = $1 and users.is_deleted = false`
	countUsersByGithubIdQuery  = `select count(*) from users where github_id = $1`
	createNewUserQuery         = `insert into users(email, name, role, department, bio, photo_url, contact, github_url, github_id, github_name) VALUES ($1, $2, $3, $4,$5, $6, $7, $8, $9, $10) returning id`
	updateUserByUserIdQuery    = `update users set email = $1, name = $2, role = $3, department = $4, bio = $5, contact = $6 where id = $7`
)
