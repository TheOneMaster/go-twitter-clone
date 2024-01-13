package db

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/TheOneMaster/go-twitter-clone/templates"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id           int
	Username     string
	DisplayName  string         `db:"displayName"`
	ProfilePhoto sql.NullString `db:"profilePhoto"`
	BannerPhoto  sql.NullString `db:"bannerPhoto"`
	CreationTime time.Time      `db:"creationTime"`
	Password     string
}

const defaultProfilePhoto = "/static/profile.png"
const defaultBannerPhoto = "/static/banner.avif"

func (user *User) VerifyExists() bool {
	var count int
	err := Connection.Get(&count, "SELECT 1 FROM Users WHERE username==?", user.Username)

	if err != nil || count == 0 {
		return false
	}

	return true
}

func (user *User) Save() error {
	hashed_password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = Connection.Exec("INSERT INTO Users(username, displayName, password) VALUES (?, ?, ?)",
		user.Username, user.DisplayName, string(hashed_password))

	if err != nil {
		slog.Error(err.Error(), "user", user)
	}

	slog.Info("Inserted user", "user", user)
	return err
}

func (user *User) ValidateLogin() bool {
	temp_store := struct {
		Count    int
		Password string
	}{}

	err := Connection.Get(&temp_store, "SELECT count(*) as count, password FROM Users WHERE username==?", user.Username)
	if err != nil {
		slog.Error(err.Error(), "user", user)
		return false
	}

	if temp_store.Count == 0 {
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(temp_store.Password), []byte(user.Password))
	return err == nil
}

func (user *User) GetDetails() error {
	err := Connection.Get(user, "SELECT username, id, displayName, profilePhoto, creationTime FROM Users WHERE username==?", user.Username)
	if err != nil {
		slog.Error(err.Error(), "user", user)
	}
	return err
}

func (user *User) GetFullDetails() error {
	err := Connection.Get(user, "SELECT * FROM Users WHERE username==?", user.Username)
	if err != nil {
		slog.Error(err.Error(), "user", user)
	}
	return err
}

func (user *User) GetImages() {
	profilePhoto := user.ProfilePhoto.String
	bannerPhoto := user.BannerPhoto.String

	if profilePhoto == "" {
		profilePhoto = defaultProfilePhoto
	}
	if bannerPhoto == "" {
		bannerPhoto = defaultBannerPhoto
	}

	user.ProfilePhoto.String = profilePhoto
	user.BannerPhoto.String = bannerPhoto
}

func (user *User) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("id", user.Id),
		slog.String("username", user.Username),
	)
}

func GetUserDetails(username string) (templates.ProfileUser, error) {
	user := User{}
	var profileDetails templates.ProfileUser
	err := Connection.Get(&user, "SELECT * FROM Users WHERE username=?", username)
	if err != nil {
		slog.Error(err.Error())
	}
	user.GetImages()

	profileDetails = templates.ProfileUser{
		Id:           user.Id,
		Username:     user.Username,
		DisplayName:  user.DisplayName,
		ProfilePhoto: user.ProfilePhoto.String,
		BannerPhoto:  user.BannerPhoto.String,
		CreationTime: user.CreationTime.Format(time.DateOnly),
	}

	return profileDetails, err
}
