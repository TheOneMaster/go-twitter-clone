package handlers

import (
	"net/http"

	"github.com/TheOneMaster/go-twitter-clone/db"
	"github.com/TheOneMaster/go-twitter-clone/templates"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	t, err := templates.LoadFiles("base.html", "login.html")
	if err != nil {
		PageNotFound(w)
	}

	err = t.Execute(w, nil)
	if err != nil {
		ServerError(w)
	}
}

type loginFormProps struct {
	Incorrect bool
}

func LoginRequest(w http.ResponseWriter, r *http.Request) {
	loginProps := loginFormProps{}

	username := r.FormValue("username")
	_ = r.FormValue("password")

	user_exists := username != "" && db.CheckUserExists(username)
	password_correct := true

	if user_exists && password_correct {
		w.Header().Set("HX-Redirect", "/")
	}

	loginProps.Incorrect = true
	t, err := templates.LoadFragment("loginForm.html")

	if err != nil {
		return
	}

	t.Execute(w, loginProps)

}
