package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const DB_FILE = "db/database.db"

var Connection *sqlx.DB

type FrontEndMessage struct {
	Author string `db:"displayName"`
	Data   string `db:"messageText"`
	Time   string
}

func DBMessages() []FrontEndMessage {
	var messages []FrontEndMessage
	rows, err := Connection.Queryx("SELECT displayName, messageText, time FROM Messages, Users WHERE author==Users.id;")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return messages
	}
	defer rows.Close()

	for rows.Next() {
		var m FrontEndMessage
		err = rows.StructScan(&m)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		messages = append(messages, m)
	}
	return messages
}

func InitConnection() {
	db, err := sqlx.Open("sqlite3", DB_FILE)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		fmt.Println("Database not properly opened")
	}

	fmt.Println("DB Connection established")
	Connection = db
}
