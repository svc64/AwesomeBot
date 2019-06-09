package main

import (
	"github.com/sonyarouje/simdb/db"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
)

type User struct {
	UID      string `json:"ID"`
	Name     string `json:"name"` // Their name
	Username string `json:"username"`
}

var driver = getDB()

//ID any struct that needs to persist should implement this function defined
//in Entity interface.
func (User User) ID() (jsonField string, value interface{}) {
	value = User.UID
	jsonField = "uid"
	return
}
func getDB() *db.Driver {
	driver, err := db.New(awesomeConfig + "/Database")
	if err != nil {
		panic(err)
	}
	return driver
}

/* This function runs every time a user sends something to check if they are in the database,
If they are not, we'll add them. */
func DBTrigger(TGUser *tb.User) {
	dbUser := getUserByID(TGUser.ID)
	if dbUser.UID == "" {
		// User doesn't exist, we are going to add them
		addUser(TGUser)
		return
	} else if dbUser.UID != strconv.Itoa(TGUser.ID) {
		panic("Fix your code, what the fuck have you done!?") // Anti-bug
	} else if dbUser.Username != TGUser.Username || dbUser.Name != getFullName(TGUser) {
		// User changed their name or username
		toDel := User{
			UID: dbUser.UID,
		}
		err := driver.Delete(toDel)
		checkGeneralError(err)
		addUser(TGUser)
	}
}
func DBWatch(b *tb.Bot) {
	b.Handle(tb.OnText, func(m *tb.Message) {
		DBTrigger(m.Sender)
		if m.IsForwarded() {
			DBTrigger(m.OriginalSender)
		}
		if m.IsReply() {
			DBTrigger(m.ReplyTo.Sender)
		}
	})
}
func addUser(User *tb.User) {
	parsedUser := parseUser(User)
	err := driver.Insert(parsedUser)
	checkGeneralError(err)
}
func parseUser(user *tb.User) User {
	fullName := getFullName(user)
	parsed := User{
		Name:     fullName,
		UID:      strconv.Itoa(user.ID),
		Username: user.Username,
	}
	return parsed
}
func getUserByName(username string) User {
	var user User
	err := driver.Open(User{}).Where("username", "=", username).First().AsEntity(&user)
	checkGeneralError(err)
	return user
}
func getUserByID(uid int) User {
	ID := strconv.Itoa(uid)
	var user User
	err := driver.Open(User{}).Where("ID", "=", ID).First().AsEntity(&user)
	checkGeneralError(err)
	return user
}
func getFullName(user *tb.User) string {
	return user.FirstName + " " + user.LastName
}
