package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/TheOneMaster/go-twitter-clone/db"
	"github.com/TheOneMaster/go-twitter-clone/templates"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimPrefix(r.URL.Path, "/profile/")
	userDetails, logged_in := isLoggedIn(r)

	if username == "me" {
		if logged_in {
			redirectURL := fmt.Sprintf("/profile/%s", userDetails.Username)
			http.Redirect(w, r, redirectURL, http.StatusFound)
		} else {
			http.Redirect(w, r, "/login", http.StatusFound)
		}
		return
	}

	editable := false
	if logged_in && username == userDetails.Username {
		editable = true
	}

	user, err := db.GetUserDetails(username)
	if err != nil {
		PageNotFound(w)
		return
	}

	messages := db.GetUserMessages(username)

	profilePageProps := templates.ProfileProps{
		User:     user,
		Messages: messages,
		Editable: editable,
	}

	t, err := templates.LoadFiles("base.html", "profile.html")
	if err != nil {
		ServerError(w)
		return
	}

	w.Header().Add("Content-Type", "text/html")

	err = t.Execute(w, profilePageProps)
	if err != nil {
		println("t1")
		slog.Error(err.Error())
	}
}
