package db

import (
	"errors"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

func ValidateLogin(username string, password string) bool {
	temp_store := struct {
		Count    int
		Password string
	}{}

	err := Connection.Get(&temp_store, "SELECT count(*) as count, password FROM Users WHERE username==?", username, password)
	if err != nil {
		slog.Error("Database call error")
		return false
	}

	if temp_store.Count == 0 {
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(temp_store.Password), []byte(password))
	return err == nil
}

type User struct {
	Username    string
	DisplayName string
	Password    string
	Photo       string
}

func SaveUser(user User) error {
	insertStatement, err := Connection.Prepare("INSERT INTO Users (username, displayName, photo, password) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}

	hashed_password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = insertStatement.Exec(user.Username, user.DisplayName, user.Photo, string(hashed_password))
	slog.Info("Insert user: ", user.Username)

	return err
}

/*
Return an error if the user does not exist and nil if the user exists
*/
func CheckUserExists(username string) error {
	var count int

	err := Connection.Get(&count, "SELECT count(*) FROM Users WHERE username==?", username)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}
	return errors.New("User does not exist")
}
