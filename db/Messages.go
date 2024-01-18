package db

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/TheOneMaster/go-twitter-clone/templates"
)

// Template for reading data from database for message
type Message struct {
	Id          int
	ParentID    sql.NullInt32 `db:"parentID"`
	MessageText string        `db:"messageText"`
	PostTime    time.Time     `db:"postTime"`
	Liked       int
	MessageUser
}

type MessageUser struct {
	Id           int `db:"author"`
	Username     string
	DisplayName  string         `db:"displayName"`
	ProfilePhoto sql.NullString `db:"profilePhoto"`
}

type Reply struct {
	Message
	Level int
}

func (msg *Message) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("id", msg.Id),
		slog.Int("author", msg.MessageUser.Id),
	)
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

func (msg *Message) GetReplies(user *User) templates.MessageList {
	msgList := templates.MessageList{}
	dbMessages := []Reply{}
	query := `
	WITH RECURSIVE reply(id, parentID, messageText, author, postTime, level) AS (
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
		SELECT r.*, displayName, profilephoto, EXISTS(
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
		tempMessage := msg.ConvertToTemplate()

		msgList = append(msgList, tempMessage)
	}

	return msgList
}

func (msg *Message) ConvertToTemplate() templates.Message {

	profilePhoto := defaultProfilePhoto
	if msg.MessageUser.ProfilePhoto.Valid {
		profilePhoto = msg.ProfilePhoto.String
	}

	msgUser := templates.MessageUser{
		Username:     msg.Username,
		DisplayName:  msg.DisplayName,
		ProfilePhoto: profilePhoto,
	}

	liked := msg.Liked == 1

	templateMessage := templates.Message{
		ID:    msg.Id,
		Data:  msg.MessageText,
		Time:  msg.PostTime.Format(time.DateTime),
		User:  msgUser,
		Liked: liked,
	}

	return templateMessage
}

func GetMessageList(user *User) templates.MessageList {
	msgList := []Message{}
	messages := templates.MessageList{}
	query := `
		SELECT M.id, messageText, postTime, displayName, username, author, profilePhoto
		FROM Messages as M, Users as U
		WHERE U.id = M.author AND parentID IS NULL
		ORDER BY postTime DESC
		LIMIT 10;
		`

	if user.VerifyExists() {
		query = `
		SELECT M.id, messageText, postTime, displayName, username, M.author, profilePhoto, EXISTS (
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
		templateMsg := msg.ConvertToTemplate()
		messages = append(messages, templateMsg)
	}

	return messages
}

func GetMessageById(msgID int, user User) (Message, error) {
	dbMsg := Message{}

	query := `
	SELECT m.id, author, displayName, profilePhoto, messageText, postTime, EXISTS (
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
	}

	return dbMsg, err
}

func GetMessagesFromUser(user *User) templates.MessageList {
	dbMsgList := []Message{}
	templateMsgList := templates.MessageList{}

	query := `
	SELECT m.id as id, u.id as author, username, messageText, postTime, displayName, profilePhoto, EXISTS(
		SELECT 1
		FROM Likes l
		WHERE l.messageid = m.id AND l.personid = u.id
	) as liked
	FROM Messages m
	RIGHT JOIN Users u ON u.id = m.author
	WHERE u.id = ?
	LIMIT 20;
	`

	err := Connection.Select(&dbMsgList, query, user.Id)
	if err != nil {
		slog.Error(err.Error())
	}

	for _, msg := range dbMsgList {
		msgTemplate := msg.ConvertToTemplate()
		templateMsgList = append(templateMsgList, msgTemplate)
	}

	return templateMsgList
}
