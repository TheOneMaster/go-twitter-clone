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
	messageID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/replies/"))
	if err != nil {
		PageNotFound(w)
		return
	}

	message, err := db.GetTopLevelMessages(messageID)

	if err != nil {
		slog.Error(err.Error(), "msg", *message)
		PageNotFound(w)
		return
	}

	user, _ := isLoggedIn(r)

	replies := message.GetReplies(user)
	println(len(replies))

	ServeFragment(w, "replyList.html", replies)
}

func LikeMessage(w http.ResponseWriter, r *http.Request) {
	user, logged_in := isLoggedIn(r)
	msgID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/like/"))

	if err != nil {
		ServerError(w)
		return
	}

	msg := db.Message{Id: msgID}

	if !logged_in || !user.VerifyExists() || !msg.VerifyExists() {
		NotAuthorizedError(w)
		return
	}

	err = msg.GetDetails()
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

	messageProps := db.GetMessage(msgID, *user)

	ServeFragment(w, "message.html", messageProps)

}
