package main
import (
	tb "gopkg.in/tucnak/telebot.v2"
	"fmt"
)
// This function is meant to run when the bot had an error and couldn't send an error message in the chat.
// err = the error
// m = the message
func handleError(err error, m tb.Message ) {
	if err != nil {
		fmt.Println("Could not send message to chat: " + m.Chat.Title)
		if m.Chat.Username != "" {
			fmt.Println("Chat username: " + m.Chat.Username)
		}
	}
}
