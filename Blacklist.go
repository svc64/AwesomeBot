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
	"bufio"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var blacklistFolder = awesomeConfig + "/blacklists/"

// handle a word blacklist
func handleBlacklist(b *tb.Bot) {
	b.Handle(tb.OnText, func(m *tb.Message) {
		words, status := readBlacklist(m.Chat.ID)
		if status { // status tells us if the chat has a blacklist or not
			lines := len(words)
			bot, err := b.ChatMemberOf(m.Chat, b.Me)
			checkError(err, m)
			var i int
			for i <= lines {
				if strings.ContainsAny(m.Text, words[i]) && bot.CanDeleteMessages {
					err = b.Delete(m)
					checkError(err, m)
					break
				}
				i++
			}
		}
	})
}

// read the blacklist file
// this function returns a string that contains the file's contents and a bool to check errors
func readBlacklist(chatID64 int64) ([]string, bool) { // chatID64 is the chat ID in int64
	if !fileExists(blacklistFolder) {
		err := os.Mkdir(blacklistFolder, 0700)
		checkGeneralError(err)
	}
	// convert int64 to string
	chatID := strconv.FormatInt(chatID64, 10)
	fileContents, err := ioutil.ReadFile(blacklistFolder + chatID) // there's a file for every blacklist
	if err != nil {
		// this error doesn't need to be reported
		return nil, false
	}
	list := string(fileContents)
	scanner := bufio.NewScanner(strings.NewReader(list))
	var words []string
	for scanner.Scan() { // count the number of lines in the config file and append them to the words array
		words = append(words, scanner.Text())
	}
	for scanner.Scan() {
	}
	return words, true
}
