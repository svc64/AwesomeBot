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
	tb "git.asafniv.me/blzit420/telebot.v2"
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
		fmt.Println("ERRRRRR")
		os.Exit(21)
	}
	b.Handle("/hello", func(m *tb.Message) {
		_, sendError := b.Send(m.Sender, "HI THERE")
		handleError(err, sendError, *m)
	})
	b.Handle("/banme", func(m *tb.Message) {
		user, _ := b.ChatMemberOf(m.Chat, m.Sender)
		if err != nil {
			_, sendError := b.Send(m.Chat, "ERRRRR")
			handleError(err, sendError, *m)
		}
		err = b.Ban(m.Chat, user)
		if err != nil {
			_, sendError := b.Send(m.Chat, "ERRRR")
			handleError(err, sendError, *m)
		}
	})
	// handle bans
	b.Handle("/ban", func(m *tb.Message) {
		sender, err := b.ChatMemberOf(m.Chat, m.Sender)
		handleError(err, nil, *m)
		bot, err := b.ChatMemberOf(m.Chat, b.Me)
		handleError(err, nil, *m)
		banUser(*b, *sender, *bot, m, false)
	})
	// kick: remove user without banning
	b.Handle("/kick", func(m *tb.Message) {
		sender, err := b.ChatMemberOf(m.Chat, m.Sender)
		handleError(err, nil, *m)
		handleError(err, nil, *m)
		bot, err := b.ChatMemberOf(m.Chat, b.Me)
		handleError(err, nil, *m)
		banUser(*b, *sender, *bot, m, true)
	})
	b.Handle("/song", func(m *tb.Message) {
		// notify the user
		status ,err := b.Reply(m, "Downloading song...")
		handleError(nil, err, *m)
		videoID, succeeded := downloadVideo(m.Payload)  // It also downloads the video and returns the ID
		err = b.Delete(status)
		handleError(err, nil, *m)
		if succeeded { // if it succeeded, send the file from disk
			sendSong(b, videoID, m)
		} else {
			status ,err = b.Reply(m, "Download failed!")
			handleError(nil, err, *m)
			// Delete after a minute
			time.Sleep(time.Minute)
			err = b.Delete(status)
		}
	})
	b.Handle("/pin", func(m *tb.Message) {
		sender, err := b.ChatMemberOf(m.Chat, m.Sender)
		handleError(err, nil, *m)
		if sender.CanPinMessages || sender.Role == tb.Creator { // check if the sender can pin messages
			err = b.Pin(m.ReplyTo)
			handleError(nil, err, *m)
		}
	})
	// id: get a user's ID
	b.Handle("/id", func(m *tb.Message) {
		if m.IsReply() { // if it's a reply we should get the ID of the user that the sender replied to
			id := strconv.Itoa(m.ReplyTo.Sender.ID) // convert an int to string
			msg := "User ID: " + id
			_ ,err = b.Send(m.Chat, msg)
			handleError(nil, err, *m)
		} else {
			id := strconv.Itoa(m.Sender.ID)
			msg := "Your ID: " + id
			_ ,err = b.Send(m.Chat, msg)
			handleError(nil, err, *m)
		}
	})
	// mid: get a message's ID
	b.Handle("/mid", func(m *tb.Message) {
		if m.IsReply() {
			id := strconv.Itoa(m.ReplyTo.ID)
			msg := "Message ID: " + id
			_ ,err = b.Reply(m, msg)
			handleError(err, nil, *m)
		} else {
			_ ,err = b.Reply(m, "Reply to a message with this command to get it's ID")
			handleError(err, nil, *m)
		}
	})
	// purge: delete every message since m.ReplyTo
	b.Handle("/purge", func(m *tb.Message) {
		sender, err := b.ChatMemberOf(m.Chat, m.Sender)
		handleError(err, nil, *m)
		bot, err := b.ChatMemberOf(m.Chat, b.Me)
		handleError(err, nil, *m)
		if sender.CanDeleteMessages ||
			tb.Creator == sender.Role && bot.CanDeleteMessages {
				err = b.Purge(m.Chat, m.ReplyTo, m)
				handleError(err, nil, *m)
		} else if !bot.CanDeleteMessages && m.Chat.Type != tb.ChatPrivate {
			_ ,err = b.Reply(m, "I don't have permission to delete messages!")
			handleError(err, nil, *m)
		} else if !sender.CanDeleteMessages {
			err = b.Delete(m)
			handleError(err, nil, *m)
		}
	})
	b.Handle("/delete", func(m *tb.Message) {
		sender, err := b.ChatMemberOf(m.Chat, m.Sender)
		handleError(err, nil, *m)
		bot, err := b.ChatMemberOf(m.Chat, b.Me)
		handleError(err, nil, *m)
		if sender.CanDeleteMessages ||
			tb.Creator == sender.Role && bot.CanDeleteMessages {
				err = b.Delete(m.ReplyTo)
				handleError(err, nil, *m)
				err = b.Delete(m)
				handleError(err, nil, *m)
		} else if !bot.CanDeleteMessages && m.Chat.Type != tb.ChatPrivate {
			_ ,err = b.Reply(m, "I don't have permission to delete messages!")
			handleError(err, nil, *m)
		} else if !sender.CanDeleteMessages && m.Chat.Type != tb.ChatPrivate && m.Sender != m.ReplyTo.Sender { // if the sender can't delete messages we'll delete their command
			err = b.Delete(m)
			handleError(err, nil, *m)
		} else if m.Chat.Type == tb.ChatPrivate && m.ReplyTo.Sender == b.Me { // the user can delete the bot's messages in PM
			err = b.Delete(m.ReplyTo)
			handleError(err, nil, *m)
			err = b.Delete(m)
			handleError(err, nil, *m)
		} else if m.Sender == m.ReplyTo.Sender && m.Chat.Type != tb.ChatPrivate { // the sender can delete their own messages
			err = b.Delete(m.ReplyTo)
			handleError(err, nil, *m)
			err = b.Delete(m)
			handleError(err, nil, *m)
		}
	})
	b.Start()
}
