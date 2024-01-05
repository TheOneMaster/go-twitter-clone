package templates

type Message struct {
	ID     int
	Author string
	Data   string
	Time   string
	Photo  string

	// Metadata properties
	Liked    bool
	Selected bool
}

type MessageList []Message

type User struct {
	Username    string
	DisplayName string
	Photo       string
}

type Profile struct {
	Username     string
	DisplayName  string
	DateCreated  string
	ProfileImage string

	Messages MessageList
}
