package handlers

import (
	"net/http"

	"github.com/TheOneMaster/go-twitter-clone/db"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	ServeStaticPage("login.html", w)
}

type loginFormProps struct {
	Incorrect bool
}

func LoginRequest(w http.ResponseWriter, r *http.Request) {
	loginProps := loginFormProps{}

	username := r.FormValue("username")
	password := r.FormValue("password")

	validate := db.ValidateLogin(username, password)

	if validate {
		logIn(w, r, username)

		redirect("/", w)
		return
	}

	loginProps.Incorrect = true
	ServeFragment(w, "loginForm.html", loginProps)
}
