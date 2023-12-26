package handlers

import (
	"fmt"
	"net/http"

	"github.com/TheOneMaster/go-twitter-clone/db"
	"github.com/TheOneMaster/go-twitter-clone/templates"
)

type IndexProps struct {
	Messages []db.FrontEndMessage
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	t, err := templates.LoadFiles("base.html", "index.html")
	if err != nil {
		fmt.Println(err)
		PageNotFound(w)
	}

	messages := db.DBMessages()

	fmt.Printf("Number of messages: %d\n", len(messages))

	pageProps := IndexProps{
		Messages: messages,
	}

	err = t.Execute(w, pageProps)
	if err != nil {
		ServerError(w)
	}
}
