package handlers

import (
	"context"
	"net/http"

	"github.com/TheOneMaster/go-twitter-clone/db"
	"github.com/TheOneMaster/go-twitter-clone/env"
	"github.com/TheOneMaster/go-twitter-clone/templates"
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

func isLoggedIn(r *http.Request) (*db.User, bool) {
	user := &db.User{}
	loggedIn := r.Context().Value(authContextKey).(bool)

	if !loggedIn {
		return user, loggedIn
	}

	session, _ := store.Get(r, sessionKey)
	username, _ := session.Values[userKey].(string)

	user.Username = username
	user.GetDetails()

	return user, loggedIn
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

func GetSidebarProps(r *http.Request) templates.SideBarProps {
	current_user, logged_in := isLoggedIn(r)
	sidebar_details := templates.SideBarProps{LoggedIn: logged_in}

	if logged_in {
		sidebar_details.User = current_user.GetSidebarDetails()
	}
	return sidebar_details
}
