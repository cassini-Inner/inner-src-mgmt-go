package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestServerUp(t *testing.T) {
	_ = os.Setenv("db_conn_string", "host=localhost port=5432 user=postgres dbname=innersource password=root sslmode=disable")

	t.Run("test mux router setup with null db", func(t *testing.T) {
		db, err := SetupRouter(nil)
		assert.NotNil(t, err, "expected err, got none")
		assert.Nil(t, db, "expected nil, got %v", db)
	})

	t.Run("setup mux router with db", func(t *testing.T) {
		db, err := sqlx.Connect("postgres", os.Getenv("db_conn_string"))
		if err != nil {
			t.Fatalf("could not connect to db: %v", err)
		}

		router, err := SetupRouter(db)
		assert.Nil(t, err, "expected nil, got: %v", err)
		assert.NotNil(t, router, "expected Mux, got nil")
	})

	t.Run("test / route", func(t *testing.T) {
		db, err:= sqlx.Connect("postgres", os.Getenv("db_conn_string"))
		if err != nil {
			t.Fatalf("could not connect to db: %v", err)
		}

		router, err := SetupRouter(db)
		assert.Nil(t, err)
		assert.NotNil(t, router)

		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("test /query route", func(t *testing.T) {
		db, err:= sqlx.Connect("postgres", os.Getenv("db_conn_string"))
		if err != nil {
			t.Fatalf("could not connect to db: %v", err)
		}

		router, err := SetupRouter(db)
		assert.Nil(t, err)
		assert.NotNil(t, router)

		request, _ := http.NewRequest(http.MethodGet, "/query", nil)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusOK, response.Code)
	})
}
