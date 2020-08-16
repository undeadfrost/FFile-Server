package middleware

import (
	"FFile-Server/db"
	"context"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func CheckLogin(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Cookie Error"))
			return
		}

		username, err := db.AuthSession(cookie.Value)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Not Login"))
			return
		}
		ctx := context.WithValue(r.Context(), "username", username)
		next(w, r.WithContext(ctx), ps)
	}
}
