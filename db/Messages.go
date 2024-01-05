package db

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/TheOneMaster/go-twitter-clone/templates"
)

type Message struct {
	Id          int
	Author      int
	ParentID    sql.NullInt32 `db:"parentID"`
	MessageText string        `db:"messageText"`
	PostTime    time.Time     `db:"postTime"`
}

func (msg *Message) Save() error {
	query := `
	INSERT INTO Messages(messageText, author)
	VALUES (?, ?);
	`
	_, err := Connection.Exec(query, msg.MessageText, msg.Author)
	slog.Info("Insert Message", "msg", msg)
	return err
}

func (msg *Message) Edit(newMessage string) error {
	_, err := Connection.Exec("UPDATE Messages SET messageText=? WHERE id==?", newMessage, msg.Id)
	return err
}

func (msg *Message) GetReplies(user *User) templates.MessageList {
	msgList := templates.MessageList{}
	dbMessages := []struct {
		Id           int
		ParentID     sql.NullInt32
		MessageText  string
		Author       string
		ProfilePhoto sql.NullString
		PostTime     time.Time
		Liked        int
		Level        int
	}{}
	query := `
	WITH RECURSIVE test(id, parentid, messagetext, author, profilephoto, posttime, level) AS (
		SELECT m.id, parentid, messagetext, username, profilephoto, posttime, 0
		  FROM Messages m, Users u
		  WHERE m.id = ? AND m.author = u.id

		  UNION ALL

		  SELECT m.id, m.parentid, m.messagetext, t.author, t.profilephoto, m.posttime, t.level+1
		  FROM Messages m
		  JOIN test t ON t.id = m.parentid
		  ORDER  BY posttime DESC
		  LIMIT 20
	)
	`
	var err error = nil
	userExists := user.VerifyExists()

	if !userExists {
		query += "SELECT * FROM test;"
		err = Connection.Select(&dbMessages, query, msg.Id)
	} else {
		query += `
		SELECT *, EXISTS(
			SELECT 1
			FROM Likes l
			WHERE l.messageID = t.id AND l.personID = ?
		) AS liked
		FROM test t;`
		err = Connection.Select(&dbMessages, query, msg.Id, user.Id)
	}

	if err != nil {
		slog.Error(err.Error(), "list", msgList)
		return msgList
	}

	for _, msg := range dbMessages {
		tempMessage := templates.Message{
			ID:       msg.Id,
			Author:   msg.Author,
			Data:     msg.MessageText,
			Time:     msg.PostTime.Format(time.DateTime),
			Photo:    msg.ProfilePhoto.String,
			Selected: false,
			Liked:    msg.Liked == 1,
		}

		msgList = append(msgList, tempMessage)
	}

	return msgList
}

func (msg *Message) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("id", msg.Id),
		slog.Int("author", msg.Author),
	)
}

func GetTopLevelMessages(id int) (*Message, error) {
	msg := &Message{}
	err := Connection.Get(msg, "SELECT * FROM Messages WHERE id==?", id)
	return msg, err
}

func GetMessageList(user *User) templates.MessageList {
	msgList := []struct {
		Id           int
		MessageText  string    `db:"messageText"`
		PostTime     time.Time `db:"postTime"`
		Author       string
		ProfilePhoto sql.NullString `db:"profilePhoto"`
		Liked        int
	}{}
	messages := templates.MessageList{}
	query := `
		SELECT M.id, messageText, postTime, username as author, profilePhoto
		FROM Messages as M, Users as U
		WHERE U.id = M.author AND parentID IS NULL
		ORDER BY postTime DESC
		LIMIT 10;
		`

	if user.VerifyExists() {
		query = `
		SELECT M.id, messageText, postTime, username as author, profilePhoto, EXISTS (
			SELECT 1
			FROM Likes AS l
			WHERE l.messageID = M.id AND l.personID = ?
		) as liked
		FROM Messages as M, Users as U
		WHERE U.id = M.author AND parentID IS NULL
		ORDER BY postTime DESC
		LIMIT 10;
		`
	}

	err := Connection.Select(&msgList, query, user.Id)

	if err != nil {
		slog.Error(err.Error())
		return messages
	}

	for _, msg := range msgList {
		temp_message := templates.Message{
			ID:     msg.Id,
			Author: msg.Author,
			Data:   msg.MessageText,
			Time:   msg.PostTime.Format(time.DateTime),
			Photo:  msg.ProfilePhoto.String,
			Liked:  msg.Liked == 1,
		}
		messages = append(messages, temp_message)
	}

	return messages
}
