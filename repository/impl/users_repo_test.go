package impl

import (
	"github.com/DATA-DOG/go-sqlmock"
	customErrors "github.com/cassini-Inner/inner-src-mgmt-go/custom_errors"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestUsersRepo_GetByIdTx(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%s was not expected while opening a stub db connection", err)
	}

	testUser := dbmodel.User{
		Id:          "1",
		Email:       dbmodel.ToNullString("test@mail.com"),
		Name:        dbmodel.ToNullString("test user name"),
		Role:        dbmodel.ToNullString("role"),
		Department:  dbmodel.ToNullString("test department"),
		Bio:         dbmodel.ToNullString("test bio"),
		PhotoUrl:    dbmodel.ToNullString("test url"),
		Contact:     dbmodel.ToNullString("test contact"),
		TimeCreated: time.Now().String(),
		TimeUpdated: time.Now().String(),
		IsDeleted:   false,
		GithubUrl:   dbmodel.ToNullString("http://github.com/test"),
		Onboarded:   false,
		GithubId:    dbmodel.ToNullString("test_id"),
		GithubName:  dbmodel.ToNullString("github_user"),
	}

	rows := sqlmock.NewRows([]string{"id", "email", "name", "role", "department", "bio", "photo_url", "contact", "time_created", "time_updated", "is_deleted", "github_url", "onboarded", "github_id", "github_name"}).
		AddRow(testUser.Id, testUser.Email, testUser.Name, testUser.Role, testUser.Department, testUser.Bio, testUser.PhotoUrl, testUser.Contact, testUser.TimeCreated, testUser.TimeUpdated, testUser.IsDeleted, testUser.GithubUrl, testUser.Onboarded, testUser.GithubId, testUser.GithubName)

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewUsersRepo(sqlxDB)
	defer db.Close()

	const query = selectUserByIdQuery

	t.Run("Pass existing id", func(t *testing.T) {
		userId := "1"

		mock.ExpectBegin()
		mock.ExpectQuery(query)
		mock.ExpectCommit()

		tx, err := sqlxDB.Beginx()
		assert.Nil(t, err)
		assert.NotNil(t, tx)

		newUser, err := repo.GetByIdTx(userId, tx)
		assert.Nil(t, err)
		assert.NotNil(t, newUser)
		err = tx.Commit()
		assert.Nil(t, err)
		if !reflect.DeepEqual(newUser, testUser) {
			t.Fatalf("got %v want %v", newUser, testUser)
		}
	})

	t.Run("Pass non-existing id", func(t *testing.T) {
		userId := "1"

		mock.ExpectBegin()
		mock.ExpectQuery("(.)*").WithArgs(userId).WillReturnError(customErrors.ErrNoEntityMatchingId)

		tx, err := sqlxDB.Beginx()

		assert.Nil(t, err)
		assert.NotNil(t, tx)
		defer tx.Commit()

		user, err := repo.GetByIdTx(userId, tx)
		assert.Nil(t, user)
		assert.NotNil(t, err)
		assert.Equal(t, err, customErrors.ErrNoEntityMatchingId)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("there are unfulfilled expectations: %s", err)
		}
	})

	t.Run("Pass invalid id", func(t *testing.T) {
		userId := "adfasdf"
		mock.ExpectBegin()
		tx, err := sqlxDB.Beginx()
		assert.Nil(t, err)
		assert.NotNil(t, tx)

		user, err := repo.GetByIdTx(userId, tx)
		assert.Nil(t, user)
		assert.NotNil(t, err)
		assert.Equal(t, err, customErrors.ErrInvalidId)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("there are unfulfilled expectations: %s", err)
		}
	})
}
