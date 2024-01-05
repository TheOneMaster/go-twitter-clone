package handlers

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/TheOneMaster/go-twitter-clone/db"
)

func SelectMessage(w http.ResponseWriter, r *http.Request) {
	messageID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/message/"))
	if err != nil {
		PageNotFound(w)
		return
	}

	message, err := db.GetTopLevelMessages(messageID)

	if err != nil {
		slog.Info(err.Error(), "msg", *message)
		PageNotFound(w)
		return
	}

	user, _ := isLoggedIn(r)

	replies := message.GetReplies(user)
	println(len(replies))

	ServeFragment(w, "replyList.html", replies)
}
