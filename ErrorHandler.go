/*
 *     AwesomeBot
 *     Copyright (C) 2019 Asaf Niv
 *
 *     This program is free software: you can redistribute it and/or modify
 *     it under the terms of the GNU General Public License as published by
 *     the Free Software Foundation, either version 3 of the License, or
 *     (at your option) any later version.
 */

package main
import (
	"github.com/getsentry/raven-go"
	tb "git.asafniv.me/blzit420/telebot.v2"
	"fmt"
	"errors"
)
//  Sentry DSN
var DSN = "https://e2672e5909514951a621c35fc7818b2d:7b34b4bae174434fa31d356f2d0d446d@sentry.io/1415174"
// This function checks if there are errors and reports them.
// err = the error
// m = the message
func handleError(err error, sendError error, m tb.Message) { // notifyChat = if the bot should sent an error message in the chat **AND FAILED TO**.
	if err != nil {
		fmt.Println(err)
		sendEvent(err, m)
	}
	if sendError != nil {
		fmt.Println(sendError)
		sendEvent(sendError, m)
	}
}
// Sends an event to sentry
func sendEvent(err error, m tb.Message) {
	sentryError := raven.SetDSN(DSN)
	if sentryError != nil {
		fmt.Println("Couldn't set a sentry DSN")
	}
	// Add details about the chat
	chatUsername := "Chat username: " + m.Chat.Username + "\n"
	chatTitle := "Chat title: " + m.Chat.Title + "\n"
	err = errors.New(chatUsername + chatTitle + err.Error())
	raven.CaptureErrorAndWait(err, nil) // Send it
}
// Handle an error that doesn't have anything to do with the chat, so there is no "m tb.Message" parameter.
func handleGeneralError(err error) {
	if err != nil {
		sentryError := raven.SetDSN(DSN)
		if sentryError != nil {
			fmt.Println("Couldn't set a sentry DSN")
		}
		raven.CaptureErrorAndWait(err, nil)
	}
}
