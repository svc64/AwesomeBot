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
	"bufio"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"strings"
)

var listFile = awesomeConfig + "/blacklist"

// handle a word blacklist
func handleBlacklist(b *tb.Bot) {
	words, status := readBlacklist()
	lines := len(words)
	if status == true {
		b.Handle(tb.OnText, func(m *tb.Message) {
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
		})
	}
}

// read the blacklist file
// this function returns a string that contains the file's contents and a bool to check errors
func readBlacklist() ([]string, bool) {
	if fileExists(listFile) {
		fileContents, err := ioutil.ReadFile(listFile)
		if err != nil {
			checkGeneralError(err)
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
	} else {
		return nil, false
	}
}
