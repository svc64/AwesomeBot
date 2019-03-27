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
	"strconv"
	"time"
)
var token string // should be provided at build time
func main() {
	b, err := tb.NewBot(tb.Settings{
		Token:  token,
		URL:    "",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		fmt.Println("Failed to set token!")
		os.Exit(21)
	}
	b.Handle("/hello", func(m *tb.Message) {
		_, err := b.Send(m.Sender, "HI THERE")
		checkError(err, m)
	})
	b.Handle("/banme", func(m *tb.Message) {
		user, err := b.ChatMemberOf(m.Chat, m.Sender)
		checkError(err, m)
		err = b.Ban(m.Chat, user)
		checkError(err, m)
	})
	// handle bans
	b.Handle("/ban", func(m *tb.Message) {
		sender, err := b.ChatMemberOf(m.Chat, m.Sender)
		checkError(err, m)
		bot, err := b.ChatMemberOf(m.Chat, b.Me)
		checkError(err, m)
		banUser(*b, *sender, *bot, m, false)
	})
	// kick: remove user without banning
	b.Handle("/kick", func(m *tb.Message) {
		sender, err := b.ChatMemberOf(m.Chat, m.Sender)
		checkError(err, m)
		bot, err := b.ChatMemberOf(m.Chat, b.Me)
		checkError(err, m)
		banUser(*b, *sender, *bot, m, true)
	})
	b.Handle("/song", func(m *tb.Message) {
		// notify the user
		status ,err := b.Reply(m, "Downloading song...")
		checkError(nil, m)
		videoID, succeeded := downloadVideo(m.Payload)  // It also downloads the video and returns the ID
		err = b.Delete(status)
		checkError(err, m)
		if succeeded { // if it succeeded, send the file from disk
			sendSong(b, videoID, m)
		} else {
			status ,err = b.Reply(m, "Download failed!")
			checkError(nil, m)
			// Delete after a minute
			time.Sleep(time.Minute)
			err = b.Delete(status)
		}
	})
	b.Handle("/pin", func(m *tb.Message) {
		sender, err := b.ChatMemberOf(m.Chat, m.Sender)
		checkError(err, m)
		if sender.CanPinMessages || sender.Role == tb.Creator { // check if the sender can pin messages
			err = b.Pin(m.ReplyTo)
			checkError(nil, m)
		}
	})
	// id: get a user's ID
	b.Handle("/id", func(m *tb.Message) {
		if m.IsReply() { // if it's a reply we should get the ID of the user that the sender replied to
			id := strconv.Itoa(m.ReplyTo.Sender.ID) // convert an int to string
			msg := "User ID: " + id
			_ ,err = b.Send(m.Chat, msg)
			checkError(err, m)
		} else {
			id := strconv.Itoa(m.Sender.ID)
			msg := "Your ID: " + id
			_ ,err = b.Send(m.Chat, msg)
			checkError(err, m)
		}
	})
	// mid: get a message's ID
	b.Handle("/mid", func(m *tb.Message) {
		if m.IsReply() {
			id := strconv.Itoa(m.ReplyTo.ID)
			msg := "Message ID: " + id
			_ ,err = b.Reply(m, msg)
			checkError(err, m)
		} else {
			_ ,err = b.Reply(m, "Reply to a message with this command to get it's ID")
			checkError(err, m)
		}
	})
	// purge: delete every message since m.ReplyTo
	b.Handle("/purge", func(m *tb.Message) {
		sender, err := b.ChatMemberOf(m.Chat, m.Sender)
		checkError(err, m)
		bot, err := b.ChatMemberOf(m.Chat, b.Me)
		checkError(err, m)
		if sender.CanDeleteMessages ||
			tb.Creator == sender.Role && bot.CanDeleteMessages {
			// The ID grows by 1 every message so we'll use a for loop and add 1 every run
			for m.ReplyTo.ID <= m.ID {
				err := b.Delete(m)
				checkError(err, m)
				m.ID++
			}
		} else if !bot.CanDeleteMessages {
			b.Reply(m, "I don't have permission to delete messages!")
		} else if !sender.CanDeleteMessages {
			err = b.Delete(m)
			checkError(err, m)
		}
	})
	b.Handle("/delete", func(m *tb.Message) {
		sender, err := b.ChatMemberOf(m.Chat, m.Sender)
		checkError(err, m)
		bot, err := b.ChatMemberOf(m.Chat, b.Me)
		checkError(err, m)
		if sender.CanDeleteMessages ||
			tb.Creator == sender.Role && bot.CanDeleteMessages {
				err = b.Delete(m.ReplyTo)
				checkError(err, m)
				err = b.Delete(m)
				checkError(err, m)
		} else if !bot.CanDeleteMessages && m.Chat.Type != tb.ChatPrivate {
			_ ,err = b.Reply(m, "I don't have permission to delete messages!")
			checkError(err, m)
		} else if !sender.CanDeleteMessages && m.Chat.Type != tb.ChatPrivate && m.Sender != m.ReplyTo.Sender { // if the sender can't delete messages we'll delete their command
			err = b.Delete(m)
			checkError(err, m)
		} else if m.Chat.Type == tb.ChatPrivate && m.ReplyTo.Sender == b.Me { // the user can delete the bot's messages in PM
			err = b.Delete(m.ReplyTo)
			checkError(err, m)
			err = b.Delete(m)
			checkError(err, m)
		} else if m.Sender == m.ReplyTo.Sender && m.Chat.Type != tb.ChatPrivate { // the sender can delete their own messages
			err = b.Delete(m.ReplyTo)
			checkError(err, m)
			err = b.Delete(m)
			checkError(err, m)
		}
	})
	b.Start()
}
