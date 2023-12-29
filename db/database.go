package db

import (
	"log"

	"github.com/TheOneMaster/go-twitter-clone/templates"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var Connection *sqlx.DB

type frontendMessage struct {
	Author string `db:"displayName"`
	Data   string `db:"messageText"`
	Time   string
	Photo  string
}

func DBMessages() templates.MessageList {
	var messages templates.MessageList
	rows, err := Connection.Queryx("SELECT displayName, messageText, time, photo FROM Messages, Users WHERE author==Users.id;")
	if err != nil {
		log.Printf("Error: %s\n", err)
		return messages
	}
	defer rows.Close()

	for rows.Next() {
		var m frontendMessage
		err = rows.StructScan(&m)
		if err != nil {
			log.Printf("Error: %s\n", err)
		}
		msg := templates.Message{
			Author: m.Author,
			Data:   m.Data,
			Photo:  m.Photo,
			Time:   m.Time,
		}
		messages = append(messages, msg)
	}
	return messages
}

func InitConnection(url string) {
	db, err := sqlx.Open("sqlite3", url)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("PRAGMA journal_mode=WAL")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("PRAGMA busy_timeout=5000")
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Println("Database not properly opened")
	}

	log.Println("DB Connection established")
	Connection = db
}
