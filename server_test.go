package main

import (
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestServerUp(t *testing.T) {
	_ = os.Setenv("db_conn_string", "host=localhost port=5432 user=postgres dbname=innersource password=root sslmode=disable")

	var setup = func() (*chi.Mux, error) {
		t.Helper()
		db, err := sqlx.Connect("postgres", os.Getenv("db_conn_string"))
		if err != nil {
			t.Fatalf("could not connect to db: %v", err)
		}

		router, err := SetupRouter(db)
		if err != nil {
			return nil, err
		}
		return router, nil
	}

	t.Run("test mux router setup with null db", func(t *testing.T) {
		db, err := SetupRouter(nil)
		assert.NotNil(t, err, "expected err, got none")
		assert.Nil(t, db, "expected nil, got %v", db)
	})

	t.Run("setup mux router with db", func(t *testing.T) {
		router, err := setup()
		assert.Nil(t, err, "expected nil, got: %v", err)
		assert.NotNil(t, router, "expected Mux, got nil")
	})

	t.Run("test / route", func(t *testing.T) {

		router, err := setup()
		assert.Nil(t, err)
		assert.NotNil(t, router)

		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("test /query route", func(t *testing.T) {
		db, err := sqlx.Connect("repository", os.Getenv("db_conn_string"))
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

	t.Run("test /authenticate ", func(t *testing.T) {
		srv, err := setup()
		assert.Nil(t, err)
		assert.NotNil(t, srv)

		request := httptest.NewRequest(http.MethodPost, "/authenticate", nil)
		response := httptest.NewRecorder()
		srv.ServeHTTP(response, request)
		assert.NotEqual(t, http.StatusNotFound, response.Code)
	})

	t.Run("test /read-cookie", func(t *testing.T) {
		srv, err := setup()
		assert.Nil(t, err)
		assert.NotNil(t, srv)

		request := httptest.NewRequest(http.MethodPost, "/read-cookie", nil)
		response := httptest.NewRecorder()
		srv.ServeHTTP(response, request)
		assert.NotEqual(t, http.StatusNotFound, response.Code)
	})

	t.Run("test /logout", func(t *testing.T) {
		srv, err := setup()
		assert.Nil(t, err)
		assert.NotNil(t, srv)

		request := httptest.NewRequest(http.MethodPost, "/logout", nil)
		response := httptest.NewRecorder()
		srv.ServeHTTP(response, request)
		assert.NotEqual(t, http.StatusNotFound, response.Code)
	})
}
