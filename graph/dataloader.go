package graph

import (
	"github.com/jmoiron/sqlx"
	"net/http"
)

func DataloaderMiddleware(db *sqlx.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
	})
}
