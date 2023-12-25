package handlers

import (
	"fmt"
	"net/http"

	"github.com/TheOneMaster/go-twitter-clone/db"
	"github.com/TheOneMaster/go-twitter-clone/templates"
)

type IndexPage struct {
	Messages []db.FrontEndMessage
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := templates.LoadFiles("layouts/base.html", "layouts/index.html", "components/username.html")
	if err != nil {
		fmt.Println(err)
		PageNotFound(w)
	}

	messages := db.GetMessages()
	pageProps := IndexPage{
		Messages: messages,
	}

	err = t.Execute(w, pageProps)
	if err != nil {
		ServerError(w)
	}

}

func PageNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404"))
}

func ServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("server error"))
}
