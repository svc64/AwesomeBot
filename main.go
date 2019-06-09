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
	tb "gopkg.in/tucnak/telebot.v2"
	"os"
	"time"
)

var token string // should be provided at build time
var youtubeAPIKey string
func main() {
	checkConfig()
	b, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		fmt.Println("Failed to set token!")
		checkError(err, nil)
		os.Exit(21)
	}
	helpMessage := readHelpMessage() // We're not reading this every time someone calls an help command.
	handleCommands(b, helpMessage)
	DBWatch(b)
	b.Start()
}
