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
import 	tb "git.asafniv.me/blzit420/telebot.v2"
// banUser: bans a user
// b: the bot
// sender: the user who called the ban command
// bot: the bot as a ChatMember
// m: the message that called the command
// kick: should we kick the user and not ban them? if kick is true, we'll unban the user right after kicking them
func banUser(b tb.Bot ,sender tb.ChatMember, bot tb.ChatMember, m *tb.Message, kick bool) {
	if sender.CanRestrictMembers || // check if the sender is an admin or the group creator.
		tb.Creator == sender.Role && bot.CanRestrictMembers { // also check if the bot can ban users
		banUser, err := b.ChatMemberOf(m.Chat, m.ReplyTo.Sender)
		handleError(err, nil, *m)
		if banUser.Role == tb.Administrator { // Check if the user is an admin or creator before banning.
			_, err := b.Reply(m, "Only the group creator can ban admins")
			handleError(nil, err, *m)
		} else if banUser.Role == tb.Creator {
			_, err := b.Reply(m, "The group creator can't be banned")
			handleError(nil, err, *m)
		} else {
			if banUser.User != bot.User { // Prevent the bot from banning itself and silently fail if someone tries to ban the bot.
				err = b.Ban(m.Chat, banUser)
				handleError(err, nil, *m)
				if err == nil && kick { // if there was no error when banning the user and kick is true, unban them after banning.
					err = b.Unban(m.Chat, banUser.User)
					handleError(err, nil, *m)
				}
			}
		}
	} else if !bot.CanRestrictMembers {
		_, sendError := b.Reply(m, "I don't have permission to ban users!")
		handleError(nil, sendError, *m)
	} else { // runs when the sender doesn't have permission to ban users
		_, sendError := b.Reply(m, "You don't have permission to ban users!")
		handleError(nil, sendError, *m)
	}
}
