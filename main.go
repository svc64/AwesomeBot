package main
import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"os"
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
		_, sendError := b.Send(m.Chat, "HI THERE")
		handleError(err, sendError, *m)
	})
	b.Handle("/banme", func(m *tb.Message) {
		user, _ := b.ChatMemberOf(m.Chat, m.Sender)
		if err != nil {
			_, sendError := b.Send(m.Chat, "ERRRRR")
			handleError(err, sendError, *m)
		}
		err := b.Ban(m.Chat, user)
		if err != nil {
			_, sendError := b.Send(m.Chat, "ERRRR")
			handleError(err, sendError, *m)
		}
	})
	// handle bans
	b.Handle("/ban", func(m *tb.Message) {
		replied := m.ReplyTo
		sender, err := b.ChatMemberOf(m.Chat, m.Sender)
		handleError(err, nil, *m)
		user, err := b.ChatMemberOf(m.Chat, m.Sender)
		handleError(err, nil, *m)
		if user.CanRestrictMembers || // check if the sender is an admin or the group creator.
			tb.Creator == sender.Role {
			user, err := b.ChatMemberOf(m.Chat, replied.Sender)
			handleError(err, nil, *m)
			err = b.Ban(m.Chat, user)
			if err != nil {
				_, sendError := b.Send(m.Chat, "ERRRR")
				handleError(err, sendError, *m)
			}
		} else { // runs when the sender doesn't have permission to ban users
			_, sendError := b.Reply(m, "You don't have permission to ban users!")
			handleError(err, sendError, *m)
		}
	})
	b.Start()
}
