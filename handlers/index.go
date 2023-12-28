package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/TheOneMaster/go-twitter-clone/db"
	"github.com/TheOneMaster/go-twitter-clone/templates"
)

type IndexProps struct {
	Messages []db.FrontEndMessage
	LoggedIn bool
	User     db.FrontEndUserDetails
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("Logged in %t\n", logIn)
	var userDetails db.FrontEndUserDetails
	username, loggedIn := isLoggedIn(r)

	if loggedIn {
		userDetails = db.GetUserDetails(username)
		log.Printf("User: %s\n", userDetails.Username)
	}

	t, err := templates.LoadFiles("base.html", "index.html")
	if err != nil {
		fmt.Println(err)
		PageNotFound(w)
	}

	messages := db.DBMessages()

	fmt.Printf("Number of messages: %d\n", len(messages))

	pageProps := IndexProps{
		Messages: messages,
		LoggedIn: loggedIn,
		User:     userDetails,
	}

	err = t.Execute(w, pageProps)
	if err != nil {
		ServerError(w)
	}
}
