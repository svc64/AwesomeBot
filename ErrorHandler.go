package main
import (
	"github.com/getsentry/raven-go"
	tb "gopkg.in/tucnak/telebot.v2"
	"fmt"
)
// This function is meant to run when the bot had an error and couldn't send an error message in the chat.
// err = the error
// m = the message
func handleError(err error, m tb.Message ) {
	if err != nil {
		sentryError := raven.SetDSN("https://e2672e5909514951a621c35fc7818b2d:7b34b4bae174434fa31d356f2d0d446d@sentry.io/1415174")
		if sentryError != nil {
			fmt.Println("Couldn't set a sentry DSN")
		}
		raven.CaptureErrorAndWait(err, nil)
		fmt.Println("Could not send message to chat: " + m.Chat.Title)
		if m.Chat.Username != "" {
			fmt.Println("Chat username: " + m.Chat.Username)
		}
	}
}
