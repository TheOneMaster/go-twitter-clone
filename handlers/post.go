package handlers

import (
	"log/slog"
	"net/http"

	"github.com/TheOneMaster/go-twitter-clone/db"
)

func MessageHandler(w http.ResponseWriter, r *http.Request) {
	_, loggedIn := isLoggedIn(r)

	if !loggedIn {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	messageText := r.FormValue("messageData")

	message := db.Message{
		MessageText: messageText,
		// Author:      user.Id,
	}

	slog.Info("new message", "message", message)

	// TODO: Properly obtain MessageText and store it in DB

	ServeFragment(w, "messageForm.html", nil)
}
