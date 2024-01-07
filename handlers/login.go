package handlers

import (
	"net/http"

	"github.com/TheOneMaster/go-twitter-clone/db"
	"github.com/TheOneMaster/go-twitter-clone/templates"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	ServeStaticPage("login.html", w)
}

func LoginRequest(w http.ResponseWriter, r *http.Request) {
	loginProps := templates.LoginFormProps{}

	username := r.FormValue("username")
	password := r.FormValue("password")

	user := db.User{
		Username: username,
		Password: password,
	}

	validate := user.ValidateLogin()
	if validate {
		logIn(w, r, username)

		redirect("/", w)
		return
	}

	loginProps.Incorrect = true
	ServeFragment(w, "loginForm.html", loginProps)
}
