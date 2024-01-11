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

	liked, err := db.ToggleLikeMessage(&msg, user)

	if err != nil {
		ServerError(w)
		return
	}

	likeButtonProps := make(map[string]any)
	likeButtonProps["ID"] = msgID
	likeButtonProps["Liked"] = liked

	reload_string := fmt.Sprintf("reload-message-%d", msgID)
	w.Header().Add("HX-Trigger", reload_string)

	ServeFragment(w, "likeButton.html", likeButtonProps)
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
