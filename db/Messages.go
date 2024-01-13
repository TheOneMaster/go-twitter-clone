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

func (msg *Message) VerifyExists() bool {
	var count int
	err := Connection.Get(&count, "SELECT 1 FROM Messages WHERE id=?", msg.Id)

	if err != nil || count == 0 {
		return false
	}

	return true
}

func (msg *Message) Save() error {
	query := `
	INSERT INTO Messages(messageText, author)
	VALUES (:messageText, :author);
	`
	_, err := Connection.Exec(query, msg)
	return err
}

func (msg *Message) Edit(newMessage string) error {
	_, err := Connection.Exec("UPDATE Messages SET messageText=? WHERE id==?", newMessage, msg.Id)
	return err
}

func (msg *Message) GetDetails() error {
	query := `
	SELECT *
	FROM Messages
	WHERE id = ?
	;`

	err := Connection.Get(msg, query, msg.Id)
	return err
}

func (msg *Message) GetReplies(user *User) templates.MessageList {
	msgList := templates.MessageList{}
	dbMessages := []struct {
		Id           int
		ParentID     sql.NullInt32
		MessageText  string
		Author       int
		DisplayName  string
		ProfilePhoto sql.NullString
		PostTime     time.Time
		Liked        int
		Level        int
	}{}
	query := `
	WITH RECURSIVE reply(id, parentid, messagetext, author, posttime, level) AS (
		SELECT m.id, parentid, messagetext, author, posttime, 0
		  FROM Messages m
		  WHERE m.id = ?

		  UNION ALL

		  SELECT m.id, m.parentid, m.messagetext, m.author, m.posttime, r.level+1
		  FROM Messages m
		  JOIN reply r ON r.id = m.parentid
		  LIMIT 20
	)
	`
	var err error = nil
	userExists := user.VerifyExists()

	if !userExists {
		query += `
		SELECT r.*, u.displayname, u.profilephoto
		FROM reply r, Users u
		WHERE r.author = u.id;
		`
		err = Connection.Select(&dbMessages, query, msg.Id)
	} else {
		query += `
		SELECT r.*, displayName as displayname, profilephoto as profilephoto, EXISTS(
			SELECT 1
			FROM Likes l
			WHERE l.messageID = r.id AND l.personID = ?
		) AS liked
		FROM reply r, Users u
		WHERE r.author = u.id;
		`
		err = Connection.Select(&dbMessages, query, msg.Id, user.Id)
	}

	if err != nil {
		slog.Error(err.Error())
		return msgList
	}

	for _, msg := range dbMessages {
		tempMessage := templates.Message{
			ID:       msg.Id,
			Author:   msg.DisplayName,
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
		Author       string
		AuthorID     string
		MessageText  string         `db:"messageText"`
		PostTime     time.Time      `db:"postTime"`
		ProfilePhoto sql.NullString `db:"profilePhoto"`
		Liked        int
	}{}
	messages := templates.MessageList{}
	query := `
		SELECT M.id, messageText, postTime, displayName as author, M.author as authorid, profilePhoto
		FROM Messages as M, Users as U
		WHERE U.id = M.author AND parentID IS NULL
		ORDER BY postTime DESC
		LIMIT 10;
		`

	if user.VerifyExists() {
		query = `
		SELECT M.id, messageText, postTime, displayName as author, M.author as authorid, profilePhoto, EXISTS (
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

func GetMessage(msgID int, user User) templates.Message {
	templateMsg := templates.Message{}
	dbMsg := struct {
		ID           int
		Author       int
		DisplayName  string         `db:"displayName"`
		ProfilePhoto sql.NullString `db:"profilePhoto"`
		MessageText  string         `db:"messageText"`
		PostTime     time.Time      `db:"postTime"`
		ParentID     int
		Liked        int
	}{}

	query := `
	SELECT m.id, author, displayName, profilePhoto, messageText, postTime, author, EXISTS (
		SELECT 1
		FROM Likes AS l
		WHERE l.messageid=m.id AND l.personid=?
	) AS liked
	FROM Messages AS m, Users as u
	WHERE m.id=? AND u.id=m.author;
	`

	err := Connection.Get(&dbMsg, query, user.Id, msgID)
	if err != nil {
		slog.Error(err.Error())
		return templateMsg
	}

	templateMsg = templates.Message{
		ID:     dbMsg.ID,
		Author: dbMsg.DisplayName,
		Data:   dbMsg.MessageText,
		Time:   dbMsg.PostTime.Format(time.DateTime),
		Photo:  dbMsg.ProfilePhoto.String,
		Liked:  dbMsg.Liked == 1,
	}

	return templateMsg
}

func GetUserMessages(username string) templates.MessageList {
	messageList := templates.MessageList{}
	dbMsgList := []messageTemplate{}

	query := `
	SELECT m.id, m.author, u.displayname, m.messagetext, m.posttime, u.profilephoto, EXISTS (
		SELECT 1
		FROM Likes l
		WHERE l.messageid = m.id AND l.personid = u.id
	) as liked
	FROM Messages m
	RIGHT JOIN Users u ON u.id = m.author
	WHERE u.username = ?;
	`

	err := Connection.Select(&dbMsgList, query, username)
	if err != nil {
		slog.Error(err.Error())
	}

	for _, msg := range dbMsgList {
		messageList = append(messageList, msg.convertToTemplateProps())
	}

	return messageList
}

type messageTemplate struct {
	Id           int
	Author       int
	ParentID     sql.NullInt32  `db:"parentID"`
	MessageText  string         `db:"messageText"`
	PostTime     time.Time      `db:"postTime"`
	DisplayName  string         `db:"displayName"`
	ProfilePhoto sql.NullString `db:"profilePhoto"`
	Liked        int
}

func (msg messageTemplate) convertToTemplateProps() templates.Message {

	profilePhoto := msg.ProfilePhoto.String
	if profilePhoto == "" {
		profilePhoto = defaultProfilePhoto
	}

	liked := msg.Liked == 1

	return templates.Message{
		ID:       msg.Id,
		Author:   msg.DisplayName,
		Data:     msg.MessageText,
		Time:     msg.PostTime.Format(time.DateTime),
		Photo:    profilePhoto,
		Liked:    liked,
		Selected: false,
	}
}
