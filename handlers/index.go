package handlers

import (
	"log/slog"
	"net/http"

	"github.com/TheOneMaster/go-twitter-clone/db"
	"github.com/TheOneMaster/go-twitter-clone/templates"
)

func IndexPage(w http.ResponseWriter, r *http.Request) {
	t, err := templates.LoadFiles("base.html", "index.html")
	if err != nil {
		slog.Error(err.Error())
		PageNotFound(w)
	}

	userDetails := templates.User{}
	user, loggedIn := isLoggedIn(r)
	if loggedIn {
		userDetails = templates.User{
			Username:    user.Username,
			DisplayName: user.DisplayName,
			Photo:       user.ProfilePhoto.String,
		}
	}

	databaseMessages := db.GetMessageList(user)

	pageProps := templates.IndexProps{
		Messages: databaseMessages,
		LoggedIn: loggedIn,
		User:     userDetails,
	}

	err = t.Execute(w, pageProps)
	if err != nil {
		ServerError(w)
	}
}
