/*
 *     AwesomeBot
 *     Copyright (C) 2019 Asaf Niv
 *
 *     This program is free software: you can redistribute it and/or modify
 *     it under the terms of the GNU Affero General Public License as published by
 *     the Free Software Foundation, either version 3 of the License, or
 *     (at your option) any later version.
 */

package main
import (
	"fmt"
	"github.com/getsentry/raven-go"
	tb "gopkg.in/tucnak/telebot.v2"
)

// DSN is the Sentry DSN
var DSN string
// This function checks if there are errors and reports them.
// err = the error
// m = the message
func checkError(err error, m *tb.Message) {
	if err != nil {
		fmt.Println(err)
		sendEvent(err, m)
	}
}
// Sends an event to sentry
func sendEvent(err error, m *tb.Message) {
	sentryError := raven.SetDSN(DSN)
	if sentryError != nil {
		fmt.Println("Couldn't set a sentry DSN")
	}
	tags := map[string]string{"Chat username: ": m.Chat.Username,
		"Chat title: ": m.Chat.Title}
	raven.CaptureErrorAndWait(err, tags) // Send it
}
// Handle an error that doesn't have anything to do with the chat, so there is no "m tb.Message" parameter here.
func checkGeneralError(err error) {
	if err != nil {
		sentryError := raven.SetDSN(DSN)
		if sentryError != nil {
			fmt.Println("Couldn't set a sentry DSN")
		}
		raven.CaptureErrorAndWait(err, nil)
	}
}
