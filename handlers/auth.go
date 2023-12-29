package handlers

import (
	"context"
	"net/http"

	"github.com/TheOneMaster/go-twitter-clone/env"
	"github.com/gorilla/sessions"
)

type ContextKey int

var store *sessions.CookieStore = nil

const sessionKey = "session"
const userKey = "userID"
const authContextKey ContextKey = 0

func ValidateSession(next http.Handler) http.Handler {
	if store == nil {
		store = sessions.NewCookieStore([]byte(env.Environment.SecretKey))
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		session, _ := store.Get(r, sessionKey)
		_, ok := session.Values[userKey]

		ctx := context.WithValue(r.Context(), authContextKey, ok)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func logIn(w http.ResponseWriter, r *http.Request, username string) {
	session, err := store.Get(r, sessionKey)
	if err != nil {
		ServerError(w)
		return
	}
	session.Values[userKey] = username
	session.Save(r, w)
}

func isLoggedIn(r *http.Request) (string, bool) {
	loggedIn := r.Context().Value(authContextKey).(bool)

	if !loggedIn {
		return "", loggedIn
	}

	session, _ := store.Get(r, sessionKey)
	username, _ := session.Values[userKey].(string)

	return username, loggedIn
}

func LogOut(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionKey)
	if err != nil {
		return
	}
	session.Values[userKey] = ""
	session.Options.MaxAge = -1
	session.Save(r, w)

	redirect("/", w)
}
