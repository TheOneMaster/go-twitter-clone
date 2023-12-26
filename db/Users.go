package db

import "fmt"

func CheckUserExists(username string) bool {
	var count int
	err := CON.Get(&count, "SELECT count(*) FROM Users WHERE username==?", username)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	return count > 0
}
