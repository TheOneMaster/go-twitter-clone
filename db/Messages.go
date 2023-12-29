package db

func SaveMessage(msg Message) error {
	statement, err := Connection.Prepare("INSERT INTO Messages(messageText, author) VALUES (?, ?)")
	if err != nil {
		return err
	}

	var userId int
	err = Connection.Get(&userId, "SELECT id FROM Users WHERE username==?", msg.Author)
	if err != nil {
		return err
	}

	_, err = statement.Exec(msg.MessageText, userId)
	return err
}
