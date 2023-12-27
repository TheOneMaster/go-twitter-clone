package handlers

import (
	"net/http"

	"github.com/TheOneMaster/go-twitter-clone/db"
	"github.com/TheOneMaster/go-twitter-clone/templates"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	ServeStaticPage("login.html", w, r)
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
		w.Header().Set("HX-Redirect", "/")
		return
	}

	loginProps.Incorrect = true
	t, err := templates.LoadFragment("loginForm.html")

	if err != nil {
		return
	}

	t.Execute(w, loginProps)

}
