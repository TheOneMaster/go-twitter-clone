package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var Connection *sqlx.DB

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
