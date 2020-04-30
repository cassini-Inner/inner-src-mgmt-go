package impl

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	customErrors "github.com/cassini-Inner/inner-src-mgmt-go/custom_errors"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/postgres/model"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestUsersRepo_GetByIdTx(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%s was not expected while opening a stub db connection", err)
	}

	testUser := dbmodel.User{
		Id:          "1",
		Email:       sql.NullString{String: "test@gmail.com"},
		Name:        sql.NullString{String: "test user"},
		Role:        sql.NullString{String: "test role"},
		Department:  sql.NullString{String: "test dept"},
		Bio:         sql.NullString{String: "test bio"},
		PhotoUrl:    sql.NullString{String: "www.github.com"},
		Contact:     sql.NullString{String: "test contact"},
		TimeCreated: "",
		TimeUpdated: "",
		IsDeleted:   false,
		GithubUrl:   sql.NullString{String: "test url"},
		Onboarded:   false,
		GithubId:    sql.NullString{String: "123456"},
		GithubName:  sql.NullString{String: "testUser"},
	}

	rows := sqlmock.NewRows([]string{"id", "email", "name", "role", "department", "bio", "photo_url", "contact", "time_created", "time_updated", "is_deleted", "github_url", "onboarded", "github_id", "github_name"}).AddRow(testUser.Id, testUser.Email, testUser.Name, testUser.Role, testUser.Department, testUser.Bio, testUser.PhotoUrl, testUser.Contact, testUser.TimeCreated, testUser.TimeUpdated, testUser.IsDeleted, testUser.GithubUrl, testUser.Onboarded, testUser.GithubId, testUser.GithubName)

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewUsersRepo(sqlxDB)
	defer db.Close()

	const query = "select"

	t.Run("Pass existing id", func(t *testing.T) {
		userId := "1223"

		mock.ExpectBegin()
		mock.ExpectQuery("(.)*").WithArgs(userId).WillReturnRows(rows)

		tx, err := sqlxDB.Beginx()

		assert.Nil(t, err)
		assert.NotNil(t, tx)
		defer tx.Commit()

		user, err := repo.GetByIdTx(userId, tx)
		assert.Nil(t, err)
		assert.NotNil(t, user)
		if !reflect.DeepEqual(user, testUser) {
			t.Fatalf("got %v want %v", user, testUser)
		}

		fmt.Println(user)
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
