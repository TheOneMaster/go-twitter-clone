package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/TheOneMaster/go-twitter-clone/db"
)

func ReplyHandler(w http.ResponseWriter, r *http.Request) {
	user, _ := isLoggedIn(r)
	messageID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/replies/"))
	if err != nil {
		PageNotFound(w)
		return
	}

	message, err := db.GetMessageById(messageID, *user)

	if err != nil {
		slog.Error(err.Error(), "msg", message)
		PageNotFound(w)
		return
	}

	replies := message.GetReplies(user)
	ServeFragment(w, "replyList.html", replies)
}

func LikeMessage(w http.ResponseWriter, r *http.Request) {
	user, logged_in := isLoggedIn(r)

	if !logged_in {
		NotAuthorizedError(w)
		return
	}

	msgID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/like/"))
	if err != nil {
		ServerError(w)
		return
	}

	msg, err := db.GetMessageById(msgID, *user)
	if err != nil {
		ServerError(w)
		return
	}

	_, err = db.ToggleLikeMessage(&msg, user)
	if err != nil {
		ServerError(w)
		return
	}

	reload_msgID := fmt.Sprintf("reload-message-%d", msgID)
	triggerEvent(reload_msgID, w)
}

func GetMessage(w http.ResponseWriter, r *http.Request) {
	msgID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/message/"))

	if err != nil {
		PageNotFound(w)
		return
	}

	user, _ := isLoggedIn(r)

	message, err := db.GetMessageById(msgID, *user)
	if err != nil {
		PageNotFound(w)
		return
	}

	messageProps := message.ConvertToTemplate()

	ServeFragment(w, "message.html", messageProps)
}
