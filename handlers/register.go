package handlers

import (
	"fmt"
	"net/http"

	"github.com/TheOneMaster/go-twitter-clone/db"
)

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	ServeStaticPage("register.html", w, r)
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

	// Currently only checking for username, no password constraints
	err := db.CheckUserExists(username)

	fmt.Printf("User Exists: %s\n", err)

	if err != nil {
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
		w.Header().Set("HX-Redirect", "/")
		return
	}

	registerProps.Username = true
	ServeFragment(w, "loginForm.html", registerProps)
}
