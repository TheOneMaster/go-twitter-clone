package templates

// Components
type Message struct {
	ID     int
	Author string
	Data   string
	Time   string
	Photo  string
	User   MessageUser

	// Metadata properties
	Liked    bool
	Selected bool
}

type MessageUser struct {
	Username     string
	DisplayName  string
	ProfilePhoto string
}

type MessageList []Message

type LoginFormProps struct {
	Incorrect bool
}

type SideBarUser struct {
	Username     string
	DisplayName  string
	ProfilePhoto string
}
type SideBarProps struct {
	LoggedIn bool
	User     SideBarUser
}

type ProfileUser struct {
	Id           int
	Username     string
	DisplayName  string
	ProfilePhoto string
	BannerPhoto  string
	CreationTime string
}

// Page Props
type IndexProps struct {
	Sidebar  SideBarProps
	Messages MessageList
}

type ProfileProps struct {
	Sidebar  SideBarProps
	User     ProfileUser
	Messages MessageList
	Editable bool
}
