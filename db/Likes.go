package db

import (
	"log/slog"
	"time"
)

type Likes struct {
	MessageID int `db:"messageID"`
	PersonID  int `db:"personID"`
	LikeTime  time.Time
}

func ToggleLikeMessage(msg *Message, user *User) (bool, error) {
	var numRows int
	Connection.Get(&numRows, "SELECT 1 FROM Likes WHERE messageID=? AND personID=?", msg.Id, user.Id)

	alreadyLiked := numRows == 1

	query := "INSERT INTO Likes(messageID, personID) VALUES (?, ?)"
	if alreadyLiked {
		query = "DELETE FROM Likes WHERE messageID=? AND personID=?"
	}

	_, err := Connection.Exec(query, msg.Id, user.Id)
	if err != nil {
		slog.Error(err.Error())
	}

	return !alreadyLiked, err
}
