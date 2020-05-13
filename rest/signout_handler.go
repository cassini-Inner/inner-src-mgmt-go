package rest

import (
	"net/http"
)

func SignoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenCookie, err := r.Cookie("token")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// setting maxage to -1 expires the cookie
	tokenCookie.MaxAge = -1
	http.SetCookie(w, tokenCookie)
}
