package templates

// Components
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

type LoginFormProps struct {
	Incorrect bool
}

// Page Props
type IndexProps struct {
	Messages MessageList
	LoggedIn bool
	User     User
}

type ProfileProps struct {
	User     ProfileUser
	Messages MessageList

	Editable bool
}

type ProfileUser struct {
	Id           int
	Username     string
	DisplayName  string
	ProfilePhoto string
	BannerPhoto  string
	CreationTime string
}
