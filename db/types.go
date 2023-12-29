package db

type Message struct {
	Id          int
	MessageText string `db:"messageText"`
	Time        string
	Author      int
}

type User struct {
	Id           int
	Username     string
	DisplayName  string `db:"displayName"`
	Photo        string
	CreationTime string `db:"creationTime"`
	Password     string
}
