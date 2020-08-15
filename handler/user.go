package handler

import (
	"FFile-Server/db"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func SignIn(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var creds = Credentials{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(creds.Username) < 5 || len(creds.Password) < 6 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	suc := db.CreateUser(creds.Username, string(hashPassword))
	if suc {
		w.Write([]byte("Success"))
	} else {
		w.Write([]byte("Failed"))
	}
}

func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var creds = Credentials{}

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	isUser := db.LoginUser(creds.Username, creds.Password)

	if !isUser {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Login Failed"))
		return
	}

	sessionToken := uuid.NewV4().String()
	isSession := db.SaveSession(creds.Username, sessionToken, 120)

	if !isSession {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ck := &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(120 * time.Second),
	}
	http.SetCookie(w, ck)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login Success"))
}
