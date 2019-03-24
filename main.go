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
		videoID, succeeded := downloadVideo(m.Payload)  // It also downloads the video and returns the ID
		if succeeded { // if it succeeded, send the file from disk
			filename := ".cache/" + videoID + ".mp4.aac"
			/* Some songs are getting an ".mp4.aac" file extension and some don't
			so we'll check for that and send a .aac file if it exists. */
			if fileExists(filename) {
				file := &tb.Audio{File: tb.FromDisk(filename)}
				_ ,err = b.Reply(m, file)
				handleError(nil, err, *m)
			} else { // song.mp4.aac doesn't exist so we'll try .aac
				filename = ".cache/" + videoID + ".aac"
				file := &tb.Audio{File: tb.FromDisk(filename)}
				_ ,err = b.Reply(m, file)
				handleError(nil, err, *m)
			}
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
			b.Reply(m, msg)
		} else {
			b.Reply(m, "Reply to a message with this command to get it's ID")
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
		} else if !bot.CanDeleteMessages {
			b.Reply(m, "I don't have permission to delete messages!")
		} else if !sender.CanDeleteMessages {
			b.Delete(m)
		}
	})
	b.Start()
}
