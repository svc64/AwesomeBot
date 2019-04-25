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
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
)
var helpMessageFile = "helpmsg"
func sendHelpMessage(b *tb.Bot, m *tb.Message, helpMessage string) {
	// Check if the chat is private and the help message isn't empty
	if m.Chat.Type == tb.ChatPrivate && helpMessage != "" {
		_, err := b.Send(m.Sender, helpMessage, tb.ModeMarkdown) // The help message is formatted in markdown
		checkError(err, m)
	}
}
// We're reading that file only once, no need to torture my disk.
func readHelpMessage() string {
	file, err := ioutil.ReadFile(awesomeConfig + "/" + helpMessageFile)
	checkGeneralError(err)
	helpMessage := string(file)
	return helpMessage
}
