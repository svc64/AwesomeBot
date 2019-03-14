package main
import (
	"github.com/getsentry/raven-go"
	tb "gopkg.in/tucnak/telebot.v2"
	"fmt"
)
// This function checks if there are errors and reports them.
// err = the error
// m = the message
func handleError(err error, sendError error, m tb.Message) { // notifyChat = if the bot should sent an error message in the chat **AND FAILED TO**.
	if err != nil {
		sendEvent(err)
	}
	if sendError != nil {
		sendEvent(sendError)
	}
}
func sendEvent(err error) {
	sentryError := raven.SetDSN("https://e2672e5909514951a621c35fc7818b2d:7b34b4bae174434fa31d356f2d0d446d@sentry.io/1415174")
	if sentryError != nil {
		fmt.Println("Couldn't set a sentry DSN")
	}
	raven.CaptureErrorAndWait(err, nil)
}
