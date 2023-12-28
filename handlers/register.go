package handlers

import (
	"net/http"

	"github.com/TheOneMaster/go-twitter-clone/db"
)

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	ServeStaticPage("register.html", w)
}

type registerFormProps struct {
	Username    bool
	ServerError bool
}

func RegisterRequest(w http.ResponseWriter, r *http.Request) {
	registerProps := registerFormProps{}

	username := r.FormValue("username")
	displayName := r.FormValue("displayName")
	password := r.FormValue("password")

	// TODO: Add password constraints
	err := db.CheckUserExists(username)

	if err == nil {
		registerProps.Username = true
		ServeFragment(w, "loginForm.html", registerProps)
		return
	}

	user := db.User{
		Username:    username,
		DisplayName: displayName,
		Password:    password,
	}
	if err = db.SaveUser(user); err != nil {
		registerProps.ServerError = true
		ServeFragment(w, "registerForm.html", registerProps)
		return
	}

	logIn(w, r, username)
	redirect("/", w)
}
