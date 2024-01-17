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

	user, _ := isLoggedIn(r)
	sidebarProps := GetSidebarProps(r)

	databaseMessages := db.GetMessageList(user)

	pageProps := templates.IndexProps{
		Sidebar:  sidebarProps,
		Messages: databaseMessages,
	}

	err = t.Execute(w, pageProps)
	if err != nil {
		ServerError(w)
	}
}
